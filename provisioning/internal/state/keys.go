package state

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID            string     `json:"id"`
	OwnerID       string     `json:"owner_id"`
	Label         string     `json:"label"` // human-redable key purpose ("internal-testing")
	HashHex       string     `json:"hash"`  // SHA256 of key
	CreatedAt     time.Time  `json:"created_at"`
	ExpiresAt     *time.Time `json:"expires_at"`     // nil - does not expire
	MaxConcurrent int        `json:"max_concurrent"` // 0 - no limit
	Revoked       bool       `json:"revoked"`
}

type keysState struct {
	Keys []APIKey `json:"keys"`
}

type KeyStore struct {
	mu    sync.Mutex
	path  string
	state keysState
}

var ErrNotFound = errors.New("key not found")

func NewKeyStore(path string) (*KeyStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	ks := &KeyStore{
		path:  path,
		state: keysState{Keys: make([]APIKey, 0)},
	}
	if err := ks.load(); err != nil {
		return nil, err
	}
	return ks, nil
}

func (ks *KeyStore) Lookup(hash [32]byte) (APIKey, bool) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	for _, k := range ks.state.Keys {
		decoded, err := hex.DecodeString(k.HashHex)
		if err != nil {
			continue
		}

		if subtle.ConstantTimeCompare(hash[:], decoded) != 1 {
			continue
		}
		if k.Revoked {
			return APIKey{}, false
		}
		if k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt) {
			return APIKey{}, false
		}
		return k, true
	}

	return APIKey{}, false
}

func (ks *KeyStore) Create(ownerID, label string, maxConcurrent int, ttl *time.Duration) (APIKey, string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return APIKey{}, "", err
	}

	b64Key := base64.URLEncoding.EncodeToString(key)
	hash := sha256.Sum256([]byte(b64Key))
	hashHex := hex.EncodeToString(hash[:])
	createdAt := time.Now().UTC()
	var expiresAt *time.Time
	if ttl != nil {
		t := createdAt.Add(*ttl)
		expiresAt = &t
	}

	apiKey := APIKey{
		ID:            uuid.NewString(),
		OwnerID:       ownerID,
		Label:         label,
		MaxConcurrent: maxConcurrent,
		Revoked:       false,
		HashHex:       hashHex,
		CreatedAt:     createdAt,
		ExpiresAt:     expiresAt,
	}

	ks.mu.Lock()
	defer ks.mu.Unlock()

	ks.state.Keys = append(ks.state.Keys, apiKey)
	if err := ks.save(); err != nil {
		ks.state.Keys = ks.state.Keys[:len(ks.state.Keys)-1]
		return APIKey{}, "", err
	}

	return apiKey, b64Key, nil
}

func (ks *KeyStore) Revoke(keyID string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	for i, k := range ks.state.Keys {
		if k.ID == keyID {
			ks.state.Keys[i].Revoked = true
			if err := ks.save(); err != nil {
				ks.state.Keys[i].Revoked = false
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("%w: unable to find key with ID %s", ErrNotFound, keyID)
}

type ExpiryOp int

const (
	ExpiryNoChange ExpiryOp = iota
	ExpirySet
	ExpiryClear
)

type ExpiryUpdate struct {
	Op    ExpiryOp
	Value *time.Time
}

func (ks *KeyStore) Update(keyID string, label *string, maxConcurrent *int, expiresAt ExpiryUpdate) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	for i, k := range ks.state.Keys {
		if k.ID == keyID {
			original := ks.state.Keys[i]

			if label != nil {
				ks.state.Keys[i].Label = *label
			}
			if maxConcurrent != nil {
				ks.state.Keys[i].MaxConcurrent = *maxConcurrent
			}
			switch expiresAt.Op {
			case ExpirySet:
				ks.state.Keys[i].ExpiresAt = expiresAt.Value
			case ExpiryClear:
				ks.state.Keys[i].ExpiresAt = nil
			}

			if err := ks.save(); err != nil {
				ks.state.Keys[i] = original
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("%w: unable to find key with ID %s", ErrNotFound, keyID)
}

func (ks *KeyStore) load() error {
	data, err := os.ReadFile(ks.path)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading keys file: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode(&ks.state)
}

func (ks *KeyStore) save() error {
	data, err := json.MarshalIndent(ks.state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	tmp := ks.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}
	return os.Rename(tmp, ks.path)
}
