package adminapi

import "time"

type createKeyRequest struct {
	OwnerID       string `json:"owner_id"`
	Label         string `json:"label"`
	MaxConcurrent int    `json:"max_concurrent"`
	TTLHours      int    `json:"ttl_hours"` // 0 - No Expiry
}

type createKeyResponse struct {
	ID        string     `json:"id"`
	Key       string     `json:"key"`
	OwnerID   string     `json:"owner_id"`
	Label     string     `json:"label"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type updateKeyRequest struct {
	Label         *string    `json:"label"`
	MaxConcurrent *int       `json:"max_concurrent"`
	ExpiresAt     *time.Time `json:"expires_at"`
	ClearExpiry   bool       `json:"clear_expiry"`
}

type getKeysResponse struct {
	Keys map[string]keyResponse `json:"keys"`
}

type keyResponse struct {
	ID            string     `json:"id"`
	OwnerID       string     `json:"owner_id"`
	Label         string     `json:"label"`
	CreatedAt     time.Time  `json:"created_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	MaxConcurrent int        `json:"max_concurrent"`
	Revoked       bool       `json:"revoked"`
}
