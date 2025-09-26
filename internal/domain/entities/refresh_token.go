package entities

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a refresh token entity in the domain
// Maps to BMSF_REFRESH_TOKEN table in Oracle database
type RefreshToken struct {
	BaseEntity
	UserID    uuid.UUID  `json:"user_id" gorm:"column:USER_ID;type:varchar(36);not null;index"` // Maps to BMSF_REFRESH_TOKEN.USER_ID
	Token     string     `json:"token" gorm:"column:TOKEN;size:500;not null;uniqueIndex"`       // Maps to BMSF_REFRESH_TOKEN.TOKEN
	ExpiresAt time.Time  `json:"expires_at" gorm:"column:EXPIRES_AT;not null;index"`            // Maps to BMSF_REFRESH_TOKEN.EXPIRES_AT
	IsRevoked bool       `json:"is_revoked" gorm:"column:IS_REVOKED;default:false;not null"`    // Maps to BMSF_REFRESH_TOKEN.IS_REVOKED
	RevokedAt *time.Time `json:"revoked_at,omitempty" gorm:"column:REVOKED_AT"`                 // Maps to BMSF_REFRESH_TOKEN.REVOKED_AT
	IPAddress string     `json:"ip_address" gorm:"column:IP_ADDRESS;size:45"`                   // Maps to BMSF_REFRESH_TOKEN.IP_ADDRESS
	UserAgent string     `json:"user_agent" gorm:"column:USER_AGENT;size:500"`                  // Maps to BMSF_REFRESH_TOKEN.USER_AGENT
}

// NewRefreshToken creates a new refresh token entity
func NewRefreshToken(userID uuid.UUID, token string, expiresAt time.Time, ipAddress, userAgent string) *RefreshToken {
	refreshToken := &RefreshToken{
		BaseEntity: NewBaseEntity(),
		UserID:     userID,
		Token:      token,
		ExpiresAt:  expiresAt,
		IsRevoked:  false,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}
	return refreshToken
}

// Revoke revokes the refresh token
func (rt *RefreshToken) Revoke(revokedBy *uuid.UUID) {
	now := time.Now()
	rt.IsRevoked = true
	rt.RevokedAt = &now
	rt.UpdateVersion(revokedBy)
}

// IsExpired checks if the refresh token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid checks if the refresh token is valid (not revoked and not expired)
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked && !rt.IsExpired()
}
