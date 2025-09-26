package entities

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity in the domain
// Maps to BMSF_USER table in Oracle database
type User struct {
	BaseEntity
	// Basic Information
	Username  string     `json:"username" gorm:"column:USERNAME;size:50;uniqueIndex;not null"`   // Maps to BMSF_USER.USERNAME
	Email     string     `json:"email" gorm:"column:EMAIL;size:255;uniqueIndex;not null"`        // Maps to BMSF_USER.EMAIL
	FirstName string     `json:"first_name" gorm:"column:FIRST_NAME;size:100;not null"`          // Maps to BMSF_USER.FIRST_NAME
	LastName  string     `json:"last_name" gorm:"column:LAST_NAME;size:100;not null"`            // Maps to BMSF_USER.LAST_NAME
	Phone     string     `json:"phone" gorm:"column:PHONE;size:20"`                              // Maps to BMSF_USER.PHONE
	Status    UserStatus `json:"status" gorm:"column:STATUS;size:20;default:'PENDING';not null"` // Maps to BMSF_USER.STATUS

	// Security Fields
	PasswordHash  string     `json:"-" gorm:"column:PASSWORD_HASH;size:255;not null"`                // Maps to BMSF_USER.PASSWORD_HASH
	Salt          string     `json:"-" gorm:"column:SALT;size:32;not null"`                          // Maps to BMSF_USER.SALT
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" gorm:"column:LAST_LOGIN_AT"`            // Maps to BMSF_USER.LAST_LOGIN_AT
	LoginAttempts int        `json:"login_attempts" gorm:"column:LOGIN_ATTEMPTS;default:0;not null"` // Maps to BMSF_USER.LOGIN_ATTEMPTS
	LockedUntil   *time.Time `json:"locked_until,omitempty" gorm:"column:LOCKED_UNTIL"`              // Maps to BMSF_USER.LOCKED_UNTIL

	// Profile Enhancement
	Avatar      string     `json:"avatar" gorm:"column:AVATAR;size:500"`                // Maps to BMSF_USER.AVATAR
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" gorm:"column:DATE_OF_BIRTH"` // Maps to BMSF_USER.DATE_OF_BIRTH
	Gender      string     `json:"gender" gorm:"column:GENDER;size:10"`                 // Maps to BMSF_USER.GENDER
	Address     string     `json:"address" gorm:"column:ADDRESS;size:500"`              // Maps to BMSF_USER.ADDRESS
	City        string     `json:"city" gorm:"column:CITY;size:100"`                    // Maps to BMSF_USER.CITY
	Country     string     `json:"country" gorm:"column:COUNTRY;size:100"`              // Maps to BMSF_USER.COUNTRY

	// Organization & Role Management
	DepartmentID *uuid.UUID `json:"department_id,omitempty" gorm:"column:DEPARTMENT_ID;type:varchar(36);index"` // Maps to BMSF_USER.DEPARTMENT_ID
	RoleID       *uuid.UUID `json:"role_id,omitempty" gorm:"column:ROLE_ID;type:varchar(36);index"`             // Maps to BMSF_USER.ROLE_ID
	ManagerID    *uuid.UUID `json:"manager_id,omitempty" gorm:"column:MANAGER_ID;type:varchar(36);index"`       // Maps to BMSF_USER.MANAGER_ID
	EmployeeCode string     `json:"employee_code" gorm:"column:EMPLOYEE_CODE;size:50;uniqueIndex"`              // Maps to BMSF_USER.EMPLOYEE_CODE

	// Notification & Preferences
	EmailVerified    bool   `json:"email_verified" gorm:"column:EMAIL_VERIFIED;default:false;not null"`               // Maps to BMSF_USER.EMAIL_VERIFIED
	PhoneVerified    bool   `json:"phone_verified" gorm:"column:PHONE_VERIFIED;default:false;not null"`               // Maps to BMSF_USER.PHONE_VERIFIED
	Language         string `json:"language" gorm:"column:LANGUAGE;size:10;default:'vi';not null"`                    // Maps to BMSF_USER.LANGUAGE
	Timezone         string `json:"timezone" gorm:"column:TIMEZONE;size:50;default:'Asia/Ho_Chi_Minh';not null"`      // Maps to BMSF_USER.TIMEZONE
	NotificationPref string `json:"notification_pref" gorm:"column:NOTIFICATION_PREF;size:20;default:'ALL';not null"` // Maps to BMSF_USER.NOTIFICATION_PREF
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
func NewUser(username, email, firstName, lastName, phone, passwordHash, salt string) *User {
	user := &User{
		BaseEntity: NewBaseEntity(),
		// Basic Information
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Status:    UserStatusPending,
		// Security Fields
		PasswordHash:  passwordHash,
		Salt:          salt,
		LoginAttempts: 0,
		// Notification & Preferences
		EmailVerified:    false,
		PhoneVerified:    false,
		Language:         "vi",
		Timezone:         "Asia/Ho_Chi_Minh",
		NotificationPref: "ALL",
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

// UpdateExtendedProfile updates extended profile information
func (u *User) UpdateExtendedProfile(avatar, gender, address, city, country string, dateOfBirth *time.Time, updatedBy *uuid.UUID) {
	u.Avatar = avatar
	u.Gender = gender
	u.Address = address
	u.City = city
	u.Country = country
	u.DateOfBirth = dateOfBirth
	u.UpdateVersion(updatedBy)
}

// UpdateOrganization updates organization information
func (u *User) UpdateOrganization(departmentID, roleID, managerID *uuid.UUID, employeeCode string, updatedBy *uuid.UUID) {
	u.DepartmentID = departmentID
	u.RoleID = roleID
	u.ManagerID = managerID
	u.EmployeeCode = employeeCode
	u.UpdateVersion(updatedBy)
}

// UpdatePreferences updates user preferences
func (u *User) UpdatePreferences(language, timezone, notificationPref string, updatedBy *uuid.UUID) {
	u.Language = language
	u.Timezone = timezone
	u.NotificationPref = notificationPref
	u.UpdateVersion(updatedBy)
}

// SetPassword updates user password
func (u *User) SetPassword(passwordHash, salt string, updatedBy *uuid.UUID) {
	u.PasswordHash = passwordHash
	u.Salt = salt
	u.UpdateVersion(updatedBy)
}

// RecordLogin records successful login
func (u *User) RecordLogin(updatedBy *uuid.UUID) {
	now := time.Now()
	u.LastLoginAt = &now
	u.LoginAttempts = 0
	u.LockedUntil = nil
	u.UpdateVersion(updatedBy)
}

// RecordFailedLogin records failed login attempt
func (u *User) RecordFailedLogin(updatedBy *uuid.UUID) {
	u.LoginAttempts++
	if u.LoginAttempts >= 5 {
		// Lock account for 30 minutes after 5 failed attempts
		lockUntil := time.Now().Add(30 * time.Minute)
		u.LockedUntil = &lockUntil
	}
	u.UpdateVersion(updatedBy)
}

// IsLocked checks if user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// UnlockAccount unlocks user account
func (u *User) UnlockAccount(updatedBy *uuid.UUID) {
	u.LoginAttempts = 0
	u.LockedUntil = nil
	u.UpdateVersion(updatedBy)
}

// VerifyEmail marks email as verified
func (u *User) VerifyEmail(updatedBy *uuid.UUID) {
	u.EmailVerified = true
	u.UpdateVersion(updatedBy)
}

// VerifyPhone marks phone as verified
func (u *User) VerifyPhone(updatedBy *uuid.UUID) {
	u.PhoneVerified = true
	u.UpdateVersion(updatedBy)
}

// GetAge calculates user age from date of birth
func (u *User) GetAge() int {
	if u.DateOfBirth == nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - u.DateOfBirth.Year()
	if now.YearDay() < u.DateOfBirth.YearDay() {
		age--
	}
	return age
}
