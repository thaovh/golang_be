package entities

import (
	"time"

	"github.com/google/uuid"
)

// BaseEntity represents the base entity with common fields
// All fields map to BMSF_* tables in Oracle database
type BaseEntity struct {
	ID        uuid.UUID  `json:"id" gorm:"column:ID;type:varchar(36);primaryKey;default:sys_guid()"` // Maps to BMSF_*.ID
	CreatedAt time.Time  `json:"created_at" gorm:"column:CREATED_AT;autoCreateTime"`                 // Maps to BMSF_*.CREATED_AT
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:UPDATED_AT;autoUpdateTime"`                 // Maps to BMSF_*.UPDATED_AT
	CreatedBy *uuid.UUID `json:"created_by,omitempty" gorm:"column:CREATED_BY;type:varchar(36)"`     // Maps to BMSF_*.CREATED_BY
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty" gorm:"column:UPDATED_BY;type:varchar(36)"`     // Maps to BMSF_*.UPDATED_BY
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:DELETED_AT;index"`                // Maps to BMSF_*.DELETED_AT
	Version   int        `json:"version" gorm:"column:VERSION;default:1;not null"`                   // Maps to BMSF_*.VERSION
	TenantID  *uuid.UUID `json:"tenant_id,omitempty" gorm:"column:TENANT_ID;type:varchar(36);index"` // Maps to BMSF_*.TENANT_ID
}

// NewBaseEntity creates a new base entity with default values
func NewBaseEntity() BaseEntity {
	return BaseEntity{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}
}

// IsDeleted checks if the entity is soft deleted
func (b *BaseEntity) IsDeleted() bool {
	return b.DeletedAt != nil
}

// SoftDelete marks the entity as deleted
func (b *BaseEntity) SoftDelete(deletedBy *uuid.UUID) {
	now := time.Now()
	b.DeletedAt = &now
	b.UpdatedAt = now
	if deletedBy != nil {
		b.UpdatedBy = deletedBy
	}
	b.Version++
}

// UpdateVersion increments the version for optimistic locking
func (b *BaseEntity) UpdateVersion(updatedBy *uuid.UUID) {
	b.UpdatedAt = time.Now()
	if updatedBy != nil {
		b.UpdatedBy = updatedBy
	}
	b.Version++
}

// Touch updates the UpdatedAt timestamp
func (b *BaseEntity) Touch(updatedBy *uuid.UUID) {
	b.UpdatedAt = time.Now()
	if updatedBy != nil {
		b.UpdatedBy = updatedBy
	}
}
