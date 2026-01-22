package domain

import (
	"time"

	"github.com/google/uuid"
)

type TrustLevel string

const (
	TrustLevelTrusted   TrustLevel = "trusted"
	TrustLevelUntrusted TrustLevel = "untrusted"
)

type Device struct {
	ID         string
	Name       string
	APIKeyHash string
	TrustLevel TrustLevel
	LastSeenAt *time.Time
	CreatedAt  time.Time
}

func NewDevice(name string, apiKeyHash string) *Device {
	return &Device{
		ID:         uuid.New().String(),
		Name:       name,
		APIKeyHash: apiKeyHash,
		TrustLevel: TrustLevelTrusted,
		CreatedAt:  time.Now(),
	}
}

func (d *Device) IsTrusted() bool {
	return d.TrustLevel == TrustLevelTrusted
}

func (d *Device) UpdateLastSeen() {
	now := time.Now()
	d.LastSeenAt = &now
}
