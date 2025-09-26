package entities

import "github.com/google/uuid"

// Department represents a department entity in the domain
// Maps to BMSF_DEPARTMENT table in Oracle database
type Department struct {
	BaseEntity
	Name        string     `json:"name" gorm:"column:NAME;size:100;not null"`                            // Maps to BMSF_DEPARTMENT.NAME
	Code        string     `json:"code" gorm:"column:CODE;size:50;uniqueIndex;not null"`                 // Maps to BMSF_DEPARTMENT.CODE
	Description string     `json:"description" gorm:"column:DESCRIPTION;size:500"`                       // Maps to BMSF_DEPARTMENT.DESCRIPTION
	ParentID    *uuid.UUID `json:"parent_id,omitempty" gorm:"column:PARENT_ID;type:varchar(36);index"`   // Maps to BMSF_DEPARTMENT.PARENT_ID
	ManagerID   *uuid.UUID `json:"manager_id,omitempty" gorm:"column:MANAGER_ID;type:varchar(36);index"` // Maps to BMSF_DEPARTMENT.MANAGER_ID
	IsActive    bool       `json:"is_active" gorm:"column:IS_ACTIVE;default:true;not null"`              // Maps to BMSF_DEPARTMENT.IS_ACTIVE
}

// NewDepartment creates a new department entity
func NewDepartment(name, code, description string, parentID, managerID *uuid.UUID) *Department {
	department := &Department{
		BaseEntity:  NewBaseEntity(),
		Name:        name,
		Code:        code,
		Description: description,
		ParentID:    parentID,
		ManagerID:   managerID,
		IsActive:    true,
	}
	return department
}

// UpdateInfo updates department information
func (d *Department) UpdateInfo(name, description string, updatedBy *uuid.UUID) {
	d.Name = name
	d.Description = description
	d.UpdateVersion(updatedBy)
}

// SetManager sets department manager
func (d *Department) SetManager(managerID *uuid.UUID, updatedBy *uuid.UUID) {
	d.ManagerID = managerID
	d.UpdateVersion(updatedBy)
}

// SetParent sets parent department
func (d *Department) SetParent(parentID *uuid.UUID, updatedBy *uuid.UUID) {
	d.ParentID = parentID
	d.UpdateVersion(updatedBy)
}

// Activate activates the department
func (d *Department) Activate(updatedBy *uuid.UUID) {
	d.IsActive = true
	d.UpdateVersion(updatedBy)
}

// Deactivate deactivates the department
func (d *Department) Deactivate(updatedBy *uuid.UUID) {
	d.IsActive = false
	d.UpdateVersion(updatedBy)
}

// IsRoot checks if this is a root department (no parent)
func (d *Department) IsRoot() bool {
	return d.ParentID == nil
}
