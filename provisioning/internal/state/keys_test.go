package state_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
	"github.com/google/uuid"
)

func defaultKeyHash() [32]byte {
	return sha256.Sum256([]byte("test-key"))
}

func defaultTestKey() state.APIKey {
	hash := defaultKeyHash()
	created := time.Now()
	return state.APIKey{
		ID:            uuid.NewString(),
		OwnerID:       "tests",
		Label:         "automated-testing",
		HashHex:       hex.EncodeToString(hash[:]),
		MaxConcurrent: 0,
		Revoked:       false,
		CreatedAt:     created,
		ExpiresAt:     nil,
	}
}

func newTestKeyStoreWithKey(t *testing.T, key state.APIKey) *state.KeyStore {
	t.Helper()
	data, err := json.Marshal(struct {
		Keys []state.APIKey `json:"keys"`
	}{Keys: []state.APIKey{key}})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "keys.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	ks, err := state.NewKeyStore(path)
	if err != nil {
		t.Fatal(err)
	}
	return ks
}

func newDefaultTestKeyStore(t *testing.T) *state.KeyStore {
	t.Helper()
	return newTestKeyStoreWithKey(t, defaultTestKey())
}

func TestLookUpKey(t *testing.T) {
	t.Run("valid-key", func(t *testing.T) {
		storedKey := defaultTestKey()
		ks := newTestKeyStoreWithKey(t, storedKey)
		key, ok := ks.Lookup(defaultKeyHash())
		if !ok {
			t.Errorf("expected to find valid key")
		}

		if key.ID != storedKey.ID {
			t.Errorf("key.ID mismatch: expected %s, got %s", storedKey.ID, key.ID)
		}
	})

	t.Run("unknown", func(t *testing.T) {
		ks := newDefaultTestKeyStore(t)
		hash := sha256.Sum256([]byte("unknown-key"))
		if _, ok := ks.Lookup(hash); ok {
			t.Errorf("expected not to find key")
		}
	})

	t.Run("revoked", func(t *testing.T) {
		key := defaultTestKey()
		key.Revoked = true
		ks := newTestKeyStoreWithKey(t, key)

		if _, ok := ks.Lookup(defaultKeyHash()); ok {
			t.Errorf("expected key to be revoked")
		}
	})

	t.Run("expired", func(t *testing.T) {
		key := defaultTestKey()
		expired := time.Now().Add(-time.Hour)
		key.ExpiresAt = &expired
		ks := newTestKeyStoreWithKey(t, key)

		if _, ok := ks.Lookup(defaultKeyHash()); ok {
			t.Errorf("expected key to be expired")
		}
	})

	t.Run("future-expiry", func(t *testing.T) {
		key := defaultTestKey()
		expired := time.Now().Add(time.Hour)
		key.ExpiresAt = &expired
		ks := newTestKeyStoreWithKey(t, key)

		if _, ok := ks.Lookup(defaultKeyHash()); !ok {
			t.Errorf("expected key to be valid")
		}
	})

	t.Run("malformed-hash", func(t *testing.T) {
		key := defaultTestKey()
		key.HashHex = "not-hex-value"
		ks := newTestKeyStoreWithKey(t, key)

		if _, ok := ks.Lookup(defaultKeyHash()); ok {
			t.Errorf("expected key to be invalid")
		}
	})
}

func TestKeysFile(t *testing.T) {
	t.Run("no-existing-keys-file", func(t *testing.T) {
		ks, err := state.NewKeyStore(filepath.Join(t.TempDir(), "state.json"))
		if err != nil {
			t.Fatalf("failed to create blank keystore")
		}

		if _, ok := ks.Lookup(defaultKeyHash()); ok {
			t.Errorf("expected no key")
		}
	})

	t.Run("invalid-json", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "keys.json")
		if err := os.WriteFile(path, []byte("not valid json"), 0600); err != nil {
			t.Fatal(err)
		}

		if _, err := state.NewKeyStore(path); err == nil {
			t.Errorf("no error thrown on broken data")
		}
	})

	t.Run("unexpected-json", func(t *testing.T) {
		data, err := json.Marshal(struct{ Foo string }{Foo: "invalid content"})
		if err != nil {
			t.Fatal(err)
		}

		path := filepath.Join(t.TempDir(), "keys.json")
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}

		if _, err := state.NewKeyStore(path); err == nil {
			t.Errorf("no error thrown on broken data")
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("valid-no-ttl", func(t *testing.T) {
		keystore := newDefaultTestKeyStore(t)
		ownerID := "new-owner-id"
		label := "internal-testing"
		maxConcurrent := 1
		k, b64Key, err := keystore.Create(ownerID, label, maxConcurrent, nil, nil)
		if err != nil {
			t.Errorf("failed to create valid key: %v", err)
		}

		if k.OwnerID != ownerID {
			t.Errorf("expected OwnerID '%s', got '%s'", ownerID, k.OwnerID)
		}

		if k.Label != label {
			t.Errorf("expected Label '%s', got '%s'", label, k.Label)
		}

		if k.MaxConcurrent != maxConcurrent {
			t.Errorf("expected MaxConcurrent '%d', got '%d'", maxConcurrent, k.MaxConcurrent)
		}

		if k.Revoked {
			t.Error("created key is revoked")
		}

		if k.ID == "" {
			t.Error("key created with blank ID")
		}

		keyHash := sha256.Sum256([]byte(b64Key))
		found, ok := keystore.Lookup(keyHash)
		if !ok {
			t.Errorf("failed to lookup created key")
		}

		if found.ID != k.ID {
			t.Errorf("key IDs do not match")
		}
	})

	t.Run("valid-with-ttl", func(t *testing.T) {
		keystore := newDefaultTestKeyStore(t)
		ownerID := "new-owner-id"
		label := "internal-testing"
		maxConcurrent := 1
		d := time.Duration(time.Hour)
		ttl := &d
		k, _, err := keystore.Create(ownerID, label, maxConcurrent, ttl, nil)
		if err != nil {
			t.Errorf("failed to create valid key: %v", err)
		}

		if k.ExpiresAt == nil {
			t.Errorf("key expected to have an expiry time")
		}

		remaining := time.Until(*k.ExpiresAt)
		if remaining < d-time.Second || remaining > d+time.Second {
			t.Errorf("key expiry time is not within tolerance")
		}
	})

	t.Run("zero-max-concurrent-rejected", func(t *testing.T) {
		keystore := newDefaultTestKeyStore(t)
		if _, _, err := keystore.Create("owner", "label", 0, nil, nil); err == nil {
			t.Error("expected error for zero max concurrent")
		}
	})

	t.Run("negative-max-concurrent-rejected", func(t *testing.T) {
		keystore := newDefaultTestKeyStore(t)
		if _, _, err := keystore.Create("owner", "label", -1, nil, nil); err == nil {
			t.Error("expected error for negative max concurrent")
		}
	})

	t.Run("persistance", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "keys.json")
		keystore, err := state.NewKeyStore(path)
		if err != nil {
			t.Fatal(err)
		}
		ownerID := "new-owner-id"
		label := "internal-testing"
		maxConcurrent := 1
		createdKey, b64Key, err := keystore.Create(ownerID, label, maxConcurrent, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		keyStore2, err := state.NewKeyStore(path)
		if err != nil {
			t.Fatal(err)
		}
		keyHash := sha256.Sum256([]byte(b64Key))
		foundKey, ok := keyStore2.Lookup(keyHash)
		if !ok {
			t.Errorf("key did not persist between keystores")
		}

		if createdKey.ID != foundKey.ID {
			t.Errorf("key IDs do not match")
		}
	})
}

func TestRevoke(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		keyStore := newDefaultTestKeyStore(t)
		k, b64Key, err := keyStore.Create("foo", "bar", 1, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := keyStore.Revoke(k.ID); err != nil {
			t.Fatal(err)
		}

		keyHash := sha256.Sum256([]byte(b64Key))
		_, ok := keyStore.Lookup(keyHash)
		if ok {
			t.Errorf("expected to fail lookup")
		}
	})

	t.Run("invalid-id", func(t *testing.T) {
		keyStore := newDefaultTestKeyStore(t)
		if err := keyStore.Revoke("some-unknown-id"); err == nil {
			t.Error("expected error on unknown id")
		}
	})

	t.Run("persistance", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "keys.json")
		keystore, err := state.NewKeyStore(path)
		if err != nil {
			t.Fatal(err)
		}

		k, b64Key, err := keystore.Create("foo", "bar", 1, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := keystore.Revoke(k.ID); err != nil {
			t.Fatal(err)
		}

		keyStore2, err := state.NewKeyStore(path)
		if err != nil {
			t.Fatal(err)
		}
		keyHash := sha256.Sum256([]byte(b64Key))
		if _, ok := keyStore2.Lookup(keyHash); ok {
			t.Errorf("revocation did not persist between keystores")
		}
	})

	t.Run("already-revoked", func(t *testing.T) {
		keyStore := newDefaultTestKeyStore(t)
		k, _, err := keyStore.Create("foo", "bar", 1, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := keyStore.Revoke(k.ID); err != nil {
			t.Fatal(err)
		}

		if err := keyStore.Revoke(k.ID); err != nil {
			t.Errorf("expected success on second call to Revoke")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("update-label", func(t *testing.T) {
		keyStore := newDefaultTestKeyStore(t)
		k, b64Key, err := keyStore.Create("foo", "original-label", 1, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		newLabel := "new-label"
		if err := keyStore.Update(k.ID, &newLabel, nil, state.ExpiryUpdate{}); err != nil {
			t.Fatal(err)
		}

		keyHash := sha256.Sum256([]byte(b64Key))
		found, ok := keyStore.Lookup(keyHash)
		if !ok {
			t.Fatal("failed to find updated key")
		}

		if found.Label != newLabel {
			t.Error("Label update did not save")
		}
	})

	t.Run("update-max-concurrent", func(t *testing.T) {
		keyStore := newDefaultTestKeyStore(t)
		k, b64Key, err := keyStore.Create("foo", "original-label", 1, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		newMaxConcurrent := 1
		if err := keyStore.Update(k.ID, nil, &newMaxConcurrent, state.ExpiryUpdate{Op: state.ExpiryNoChange}); err != nil {
			t.Fatal(err)
		}

		keyHash := sha256.Sum256([]byte(b64Key))
		found, ok := keyStore.Lookup(keyHash)
		if !ok {
			t.Fatal("failed to find updated key")
		}

		if found.MaxConcurrent != newMaxConcurrent {
			t.Error("MaxConcurrent update did not save")
		}
	})

	t.Run("zero-max-concurrent-rejected", func(t *testing.T) {
		keyStore := newDefaultTestKeyStore(t)
		k, b64Key, err := keyStore.Create("foo", "original-label", 1, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		zero := 0
		if err := keyStore.Update(k.ID, nil, &zero, state.ExpiryUpdate{Op: state.ExpiryNoChange}); err == nil {
			t.Error("expected error for zero max concurrent")
		}

		keyHash := sha256.Sum256([]byte(b64Key))
		found, ok := keyStore.Lookup(keyHash)
		if !ok {
			t.Fatal("failed to find key")
		}
		if found.MaxConcurrent != 1 {
			t.Error("rejected update should not change stored MaxConcurrent")
		}
	})

	t.Run("set-expiry", func(t *testing.T) {
		keyStore := newDefaultTestKeyStore(t)
		k, b64Key, err := keyStore.Create("foo", "original-label", 1, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		newExpiry := k.CreatedAt.Add(time.Hour)
		if err := keyStore.Update(k.ID, nil, nil, state.ExpiryUpdate{
			Op:    state.ExpirySet,
			Value: &newExpiry,
		}); err != nil {
			t.Fatal(err)
		}

		keyHash := sha256.Sum256([]byte(b64Key))
		found, ok := keyStore.Lookup(keyHash)
		if !ok {
			t.Fatal("failed to find updated key")
		}

		if found.ExpiresAt == nil {
			t.Fatal("ExpiresAt was not set")
		}

		if !found.ExpiresAt.Equal(newExpiry) {
			t.Error("Expiry time update did not save")
		}
	})

	t.Run("clear-expiry", func(t *testing.T) {
		keyStore := newDefaultTestKeyStore(t)
		ttl := time.Hour
		k, b64Key, err := keyStore.Create("foo", "original-label", 1, &ttl, nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := keyStore.Update(k.ID, nil, nil, state.ExpiryUpdate{Op: state.ExpiryClear}); err != nil {
			t.Fatal(err)
		}

		keyHash := sha256.Sum256([]byte(b64Key))
		found, ok := keyStore.Lookup(keyHash)
		if !ok {
			t.Fatal("failed to find updated key")
		}

		if found.ExpiresAt != nil {
			t.Error("Expiry time update did not clear")
		}
	})

	t.Run("invalid-id", func(t *testing.T) {
		keyStore := newDefaultTestKeyStore(t)
		if err := keyStore.Update("some-invalid-id", nil, nil, state.ExpiryUpdate{}); err == nil {
			t.Error("expected error for invalid id")
		}
	})
}

func TestReload(t *testing.T) {
	t.Run("valid-key-added", func(t *testing.T) {
		data, err := json.Marshal(struct {
			Keys []state.APIKey `json:"keys"`
		}{Keys: []state.APIKey{defaultTestKey()}})
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(t.TempDir(), "keys.json")
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		ks, err := state.NewKeyStore(path)
		if err != nil {
			t.Fatal(err)
		}

		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var keys struct {
			Keys []state.APIKey `json:"keys"`
		}
		if err := json.Unmarshal(b, &keys); err != nil {
			t.Fatal(err)
		}
		newKeyHash := sha256.Sum256([]byte("new-key"))
		newKey := state.APIKey{
			ID:      uuid.NewString(),
			HashHex: hex.EncodeToString(newKeyHash[:]),
		}
		keys.Keys = append(keys.Keys, newKey)
		d, err := json.Marshal(keys)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, d, 0600); err != nil {
			t.Fatal(err)
		}

		if _, ok := ks.Lookup(newKeyHash); !ok {
			t.Error("failed to lookup key added by admin")
		}
	})

	t.Run("unexpected-json-saved", func(t *testing.T) {
		data, err := json.Marshal(struct {
			Keys []state.APIKey `json:"keys"`
		}{Keys: []state.APIKey{defaultTestKey()}})
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(t.TempDir(), "keys.json")
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		ks, err := state.NewKeyStore(path)
		if err != nil {
			t.Fatal(err)
		}

		var brokenStruct struct {
			Keys []state.APIKey `json:"broken"`
		}
		d, err := json.Marshal(brokenStruct)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, d, 0600); err != nil {
			t.Fatal(err)
		}

		if _, ok := ks.Lookup(defaultKeyHash()); !ok {
			t.Error("failed to lookup existing valid key in memory")
		}
	})
}
