package entities

import "github.com/google/uuid"

// Role represents a role entity in the domain
// Maps to BMSF_ROLE table in Oracle database
type Role struct {
	BaseEntity
	Name        string `json:"name" gorm:"column:NAME;size:100;not null"`                // Maps to BMSF_ROLE.NAME
	Code        string `json:"code" gorm:"column:CODE;size:50;uniqueIndex;not null"`     // Maps to BMSF_ROLE.CODE
	Description string `json:"description" gorm:"column:DESCRIPTION;size:500"`           // Maps to BMSF_ROLE.DESCRIPTION
	Permissions string `json:"permissions" gorm:"column:PERMISSIONS;type:CLOB"`          // Maps to BMSF_ROLE.PERMISSIONS (JSON)
	IsActive    bool   `json:"is_active" gorm:"column:IS_ACTIVE;default:true;not null"`  // Maps to BMSF_ROLE.IS_ACTIVE
	IsSystem    bool   `json:"is_system" gorm:"column:IS_SYSTEM;default:false;not null"` // Maps to BMSF_ROLE.IS_SYSTEM
}

// NewRole creates a new role entity
func NewRole(name, code, description, permissions string, isSystem bool) *Role {
	role := &Role{
		BaseEntity:  NewBaseEntity(),
		Name:        name,
		Code:        code,
		Description: description,
		Permissions: permissions,
		IsActive:    true,
		IsSystem:    isSystem,
	}
	return role
}

// UpdateInfo updates role information
func (r *Role) UpdateInfo(name, description string, updatedBy *uuid.UUID) {
	r.Name = name
	r.Description = description
	r.UpdateVersion(updatedBy)
}

// UpdatePermissions updates role permissions
func (r *Role) UpdatePermissions(permissions string, updatedBy *uuid.UUID) {
	r.Permissions = permissions
	r.UpdateVersion(updatedBy)
}

// Activate activates the role
func (r *Role) Activate(updatedBy *uuid.UUID) {
	r.IsActive = true
	r.UpdateVersion(updatedBy)
}

// Deactivate deactivates the role
func (r *Role) Deactivate(updatedBy *uuid.UUID) {
	r.IsActive = false
	r.UpdateVersion(updatedBy)
}

// IsSystemRole checks if this is a system role
func (r *Role) IsSystemRole() bool {
	return r.IsSystem
}
