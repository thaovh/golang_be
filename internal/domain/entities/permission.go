package entities

import "github.com/google/uuid"

// Permission represents a permission entity in the domain
// Maps to BMSF_PERMISSION table in Oracle database
type Permission struct {
	BaseEntity
	Name        string `json:"name" gorm:"column:NAME;size:100;not null"`               // Maps to BMSF_PERMISSION.NAME
	Code        string `json:"code" gorm:"column:CODE;size:50;uniqueIndex;not null"`    // Maps to BMSF_PERMISSION.CODE
	Resource    string `json:"resource" gorm:"column:RESOURCE;size:100;not null"`       // Maps to BMSF_PERMISSION.RESOURCE
	Action      string `json:"action" gorm:"column:ACTION;size:50;not null"`            // Maps to BMSF_PERMISSION.ACTION
	Description string `json:"description" gorm:"column:DESCRIPTION;size:500"`          // Maps to BMSF_PERMISSION.DESCRIPTION
	IsActive    bool   `json:"is_active" gorm:"column:IS_ACTIVE;default:true;not null"` // Maps to BMSF_PERMISSION.IS_ACTIVE
}

// NewPermission creates a new permission entity
func NewPermission(name, code, resource, action, description string) *Permission {
	permission := &Permission{
		BaseEntity:  NewBaseEntity(),
		Name:        name,
		Code:        code,
		Resource:    resource,
		Action:      action,
		Description: description,
		IsActive:    true,
	}
	return permission
}

// UpdateInfo updates permission information
func (p *Permission) UpdateInfo(name, description string, updatedBy *uuid.UUID) {
	p.Name = name
	p.Description = description
	p.UpdateVersion(updatedBy)
}

// Activate activates the permission
func (p *Permission) Activate(updatedBy *uuid.UUID) {
	p.IsActive = true
	p.UpdateVersion(updatedBy)
}

// Deactivate deactivates the permission
func (p *Permission) Deactivate(updatedBy *uuid.UUID) {
	p.IsActive = false
	p.UpdateVersion(updatedBy)
}

// GetFullCode returns the full permission code (resource:action)
func (p *Permission) GetFullCode() string {
	return p.Resource + ":" + p.Action
}
