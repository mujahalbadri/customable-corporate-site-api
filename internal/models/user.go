package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"unique;not null;index"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Role      string         `json:"role" gorm:"default:'user';index"`
	IsActive  bool           `json:"is_active" gorm:"default:true;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// User roles constants
const (
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleUser   = "user"
)

// BeforeCreate hook to hash password before saving
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Hash the password if it's not empty
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	// Set default role if not set
	if u.Role == "" {
		u.Role = RoleUser
	}

	return nil
}

// BeforeUpdate hook to hash password if it has changed
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	// Only hash the password if it has been changed
	if tx.Statement.Changed("Password") && u.Password != "" {
		// Check if it's already hashed (bcrypt hashed passwords start with $2a$, $2b$, $2x$ or $2y$)
		if len(u.Password) < 60 || (u.Password[:4] != "$2a$" && u.Password[:4] != "$2b$" && u.Password[:4] != "$2x$" && u.Password[:4] != "$2y$") {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			u.Password = string(hashedPassword)
		}
	}
	return nil
}

// CheckPassword verifies the provided password against the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsEditor checks if the user has editor role
func (u *User) IsEditor() bool {
	return u.Role == RoleEditor
}

// IsUser checks if the user has user role
func (u *User) IsUser() bool {
	return u.Role == RoleUser
}

// TableName sets the insert table name for this struct type
func (User) TableName() string {
	return "users"
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		FullName:  u.GetFullName(),
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
