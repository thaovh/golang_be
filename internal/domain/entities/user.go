package entities

import "github.com/google/uuid"

// User represents a user entity in the domain
// Maps to BMSF_USER table in Oracle database
type User struct {
	BaseEntity
	Username  string     `json:"username" gorm:"column:USERNAME;size:50;uniqueIndex;not null"`   // Maps to BMSF_USER.USERNAME
	Email     string     `json:"email" gorm:"column:EMAIL;size:255;uniqueIndex;not null"`        // Maps to BMSF_USER.EMAIL
	FirstName string     `json:"first_name" gorm:"column:FIRST_NAME;size:100;not null"`          // Maps to BMSF_USER.FIRST_NAME
	LastName  string     `json:"last_name" gorm:"column:LAST_NAME;size:100;not null"`            // Maps to BMSF_USER.LAST_NAME
	Phone     string     `json:"phone" gorm:"column:PHONE;size:20"`                              // Maps to BMSF_USER.PHONE
	Status    UserStatus `json:"status" gorm:"column:STATUS;size:20;default:'PENDING';not null"` // Maps to BMSF_USER.STATUS
}

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusInactive UserStatus = "INACTIVE"
	UserStatusPending  UserStatus = "PENDING"
	UserStatusBlocked  UserStatus = "BLOCKED"
)

// IsValid checks if the user status is valid
func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusInactive, UserStatusPending, UserStatusBlocked:
		return true
	default:
		return false
	}
}

// NewUser creates a new user entity
func NewUser(username, email, firstName, lastName, phone string) *User {
	user := &User{
		BaseEntity: NewBaseEntity(),
		Username:   username,
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		Phone:      phone,
		Status:     UserStatusPending,
	}
	return user
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// Activate activates the user
func (u *User) Activate(updatedBy *uuid.UUID) {
	u.Status = UserStatusActive
	u.UpdateVersion(updatedBy)
}

// Deactivate deactivates the user
func (u *User) Deactivate(updatedBy *uuid.UUID) {
	u.Status = UserStatusInactive
	u.UpdateVersion(updatedBy)
}

// Block blocks the user
func (u *User) Block(updatedBy *uuid.UUID) {
	u.Status = UserStatusBlocked
	u.UpdateVersion(updatedBy)
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(firstName, lastName, phone string, updatedBy *uuid.UUID) {
	u.FirstName = firstName
	u.LastName = lastName
	u.Phone = phone
	u.UpdateVersion(updatedBy)
}
