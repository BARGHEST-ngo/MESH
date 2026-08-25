package state

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID            string         `json:"id"`
	OwnerID       string         `json:"owner_id"`
	Label         string         `json:"label"` // human-redable key purpose ("internal-testing")
	HashHex       string         `json:"hash"`  // SHA256 of key
	CreatedAt     time.Time      `json:"created_at"`
	ExpiresAt     *time.Time     `json:"expires_at"` // nil - does not expire
	MaxConcurrent int            `json:"max_concurrent"`
	Revoked       bool           `json:"revoked"`
	DeploymentTTL *time.Duration `json:"deployment_ttl"` // nil - use default
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

	if err := ks.reload(); err != nil {
		// log error, lookup keys in memory
	}

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

func (ks *KeyStore) Create(ownerID, label string, maxConcurrent int, keyTTL, deploymentTTL *time.Duration) (APIKey, string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return APIKey{}, "", err
	}

	if maxConcurrent <= 0 {
		return APIKey{}, "", fmt.Errorf("maxConcurrent must be greater than 0")
	}

	b64Key := base64.URLEncoding.EncodeToString(key)
	hash := sha256.Sum256([]byte(b64Key))
	hashHex := hex.EncodeToString(hash[:])
	createdAt := time.Now().UTC()
	var expiresAt *time.Time
	if keyTTL != nil {
		t := createdAt.Add(*keyTTL)
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
		DeploymentTTL: deploymentTTL,
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
				if *maxConcurrent <= 0 {
					ks.state.Keys[i] = original
					return fmt.Errorf("maxConcurrent must be greater than 0")
				}
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

func (ks *KeyStore) List() []APIKey {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	if err := ks.reload(); err != nil {
		// log error but list keys in memory
	}
	out := make([]APIKey, len(ks.state.Keys))
	copy(out, ks.state.Keys)
	return out
}

func (ks *KeyStore) load() error {
	return loadJSON(ks.path, &ks.state)
}

func (ks *KeyStore) save() error {
	return saveJSON(ks.path, ks.state)
}

func (ks *KeyStore) reload() error {
	loaded := keysState{}
	if err := loadJSON(ks.path, &loaded); err != nil {
		return err
	}
	ks.state = loaded
	return nil
}
