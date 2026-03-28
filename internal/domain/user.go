package domain

import (
	"errors"
	"time"
)

// Role defines a user role within the system.
type Role string

const (
	RoleCustomer Role = "customer"
	RoleAdmin    Role = "admin"
)

var (
	ErrInvalidRole         = errors.New("invalid user role")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("amount must be positive")
	ErrUserNotFound        = errors.New("user not found")
)

// User represents an application user.
//
// A user may authenticate either via Telegram or email/password.
// Admin user require explicit access expiration configuration.
type User struct {
	BaseAggregate

	id           int
	tgID         *int64
	tgName       *string
	email        string
	passwordHash string
	role         Role
	balance      int64
	isEnabled    bool

	adminAccessExpiresAt *time.Time
	createdAt            time.Time
	updatedAt            time.Time
}

type NewUserParams struct {
	TgID         *int64
	TgName       string
	Email        string
	PasswordHash string
	Role         Role
}

// NewUser created a new user instance.
//
// Business rules:
//   - If role is empty, it defaults to RoleCustomer.
//   - Role must be either RoleCustomer or RoleAdmin.
//   - User must authenticate either via Telegram (TgID)
//     or via email + password hash.
//
// Admin users are enabled by default.
// Customer users are disabled by default.
func NewUser(p NewUserParams) (*User, error) {
	now := time.Now()

	if p.Role == "" {
		p.Role = RoleCustomer
	}

	if p.Role != RoleCustomer && p.Role != RoleAdmin {
		return nil, ErrInvalidRole
	}

	if p.TgID == nil && (p.Email == "" || p.PasswordHash == "") {
		return nil, ErrInvalidCredentials
	}

	var tgName *string
	if p.TgID != nil {
		tgName = &p.TgName
	}

	user := &User{
		tgID:         p.TgID,
		tgName:       tgName,
		email:        p.Email,
		role:         p.Role,
		passwordHash: p.PasswordHash,
		isEnabled:    false,
		balance:      0,
		createdAt:    now,
		updatedAt:    now,
	}

	user.setInitialVersion(1)

	if user.role == RoleAdmin {
		user.isEnabled = true
	}

	return user, nil
}

// ---- GETTERS ----

// ID returns user id.
func (u *User) ID() int {
	return u.id
}

// TelegramID returns user telegram id.
func (u *User) TelegramID() (int64, bool) {
	if u.tgID == nil {
		return 0, false
	}
	return *u.tgID, true
}

// TelegramName returns user telegram name.
func (u *User) TelegramName() (string, bool) {
	if u.tgName == nil {
		return "", false
	}
	return *u.tgName, true
}

func (u *User) Balance() int64 {
	return u.balance
}

// Version returns version of user application.
func (u *User) Version() int {
	return u.version
}

// Email returns email of user.
func (u *User) Email() string {
	return u.email
}

// PasswordHash returns password hash of user.
func (u *User) PasswordHash() string {
	return u.passwordHash
}

// Role returns role of user.
func (u *User) Role() Role {
	return u.role
}

// IsEnabled returns ....
func (u *User) IsEnabled() bool {
	return u.isEnabled
}

// AdminAccessExpiresAt return time where admin access is expired.
func (u *User) AdminAccessExpiresAt() *time.Time {
	return u.adminAccessExpiresAt
}

// AddBalance increase user balance.
func (u *User) AddBalance(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	u.balance += amount
	u.updatedAt = time.Now()

	return nil
}

// DeductBalance decrease user balance.
//
// Fails if balance is insufficient.
func (u *User) DeductBalance(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if u.balance < amount {
		return ErrInsufficientBalance
	}

	u.balance -= amount
	u.updatedAt = time.Now()

	return nil
}

// CanUseAdminPanel determines whether the user
// has currently valid admin panel access.
//
// Conditions:
//   - Role must be RoleAdmin
//   - User must be enabled
//   - Access expiration must be set
//   - Current time must be before expiration
func (u *User) CanUseAdminPanel(now time.Time) bool {
	if u.role != RoleAdmin || !u.isEnabled {
		return false
	}

	if u.adminAccessExpiresAt == nil {
		return false
	}

	return now.Before(*u.adminAccessExpiresAt)
}

// ---- SETTERS ----

// Enable activates the user account
// and updates the modification timestamp.
func (u *User) Enable() {
	u.isEnabled = true
	u.updatedAt = time.Now()
}

// Disable deactivates the user account
// and updates the modification timestamp.
func (u *User) Disable() {
	u.isEnabled = false
	u.updatedAt = time.Now()
}

// GrantAdminAccess grants temporary admin access
// until the specified time.
//
// Has effect only if user role is RoleAdmin.
func (u *User) GrantAdminAccess(until time.Time) {
	if u.role != RoleAdmin {
		return
	}
	u.adminAccessExpiresAt = &until
	u.updatedAt = time.Now()
}

// ---- SETTERS ----

// SetID is intended for repository layer only.
func (u *User) SetID(id int) {
	u.id = id
}
