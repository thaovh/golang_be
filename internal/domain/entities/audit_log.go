package entities

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an audit log entity in the domain
// Maps to BMSF_AUDIT_LOG table in Oracle database
type AuditLog struct {
	BaseEntity
	UserID     *uuid.UUID `json:"user_id,omitempty" gorm:"column:USER_ID;type:varchar(36);index"`         // Maps to BMSF_AUDIT_LOG.USER_ID
	Action     string     `json:"action" gorm:"column:ACTION;size:100;not null"`                          // Maps to BMSF_AUDIT_LOG.ACTION
	Resource   string     `json:"resource" gorm:"column:RESOURCE;size:100;not null"`                      // Maps to BMSF_AUDIT_LOG.RESOURCE
	ResourceID *uuid.UUID `json:"resource_id,omitempty" gorm:"column:RESOURCE_ID;type:varchar(36);index"` // Maps to BMSF_AUDIT_LOG.RESOURCE_ID
	OldValues  string     `json:"old_values" gorm:"column:OLD_VALUES;type:CLOB"`                          // Maps to BMSF_AUDIT_LOG.OLD_VALUES (JSON)
	NewValues  string     `json:"new_values" gorm:"column:NEW_VALUES;type:CLOB"`                          // Maps to BMSF_AUDIT_LOG.NEW_VALUES (JSON)
	IPAddress  string     `json:"ip_address" gorm:"column:IP_ADDRESS;size:45"`                            // Maps to BMSF_AUDIT_LOG.IP_ADDRESS
	UserAgent  string     `json:"user_agent" gorm:"column:USER_AGENT;size:500"`                           // Maps to BMSF_AUDIT_LOG.USER_AGENT
	SessionID  string     `json:"session_id" gorm:"column:SESSION_ID;size:100"`                           // Maps to BMSF_AUDIT_LOG.SESSION_ID
	Timestamp  time.Time  `json:"timestamp" gorm:"column:TIMESTAMP;autoCreateTime;not null"`              // Maps to BMSF_AUDIT_LOG.TIMESTAMP
}

// NewAuditLog creates a new audit log entity
func NewAuditLog(userID *uuid.UUID, action, resource string, resourceID *uuid.UUID, oldValues, newValues, ipAddress, userAgent, sessionID string) *AuditLog {
	auditLog := &AuditLog{
		BaseEntity: NewBaseEntity(),
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		OldValues:  oldValues,
		NewValues:  newValues,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		SessionID:  sessionID,
		Timestamp:  time.Now(),
	}
	return auditLog
}

// GetActionDescription returns a human-readable action description
func (a *AuditLog) GetActionDescription() string {
	switch a.Action {
	case "CREATE":
		return "Created " + a.Resource
	case "UPDATE":
		return "Updated " + a.Resource
	case "DELETE":
		return "Deleted " + a.Resource
	case "LOGIN":
		return "User logged in"
	case "LOGOUT":
		return "User logged out"
	case "VIEW":
		return "Viewed " + a.Resource
	default:
		return a.Action + " " + a.Resource
	}
}
