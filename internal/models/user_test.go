package models

import (
	"testing"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the User model
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("Failed to auto-migrate test database: %v", err)
	}

	return db
}

func TestUserPasswordHashing(t *testing.T) {
	db := setupTestDB(t)

	user := &User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Role:      RoleUser,
	}

	// Create user (should trigger BeforeCreate hook)
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user in test database: %v", err)
	}

	// Password should be hashed
	if user.Password == "password123" {
		t.Errorf("Password was not hashed")
	}

	// Should be bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123")); err != nil {
		t.Errorf("Password hash does not match original password: %v", err)
	}
}

func TestUserCheckPassword(t *testing.T) {
	db := setupTestDB(t)

	user := &User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Role:      RoleUser,
	}

	// Create user (should trigger BeforeCreate hook)
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user in test database: %v", err)
	}

	// Test password verification
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"CorrectPassword", "password123", true},
		{"IncorrectPassword", "wrongpassword", false},
		{"EmptyPassword", "", false},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := user.CheckPassword(tt.password)
			if got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserGetFullName(t *testing.T) {
	user := &User{
		FirstName: "John",
		LastName:  "Doe",
	}

	want := "John Doe"
	if got := user.GetFullName(); got != want {
		t.Errorf("GetFullName() = %v, want %v", got, want)
	}
}

func TestUserRoleChecks(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		isAdmin  bool
		isEditor bool
		isUser   bool
	}{
		{"Admin", RoleAdmin, true, false, false},
		{"Editor", RoleEditor, false, true, false},
		{"User", RoleUser, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			if got := user.IsAdmin(); got != tt.isAdmin {
				t.Errorf("IsAdmin() = %v, want %v", got, tt.isAdmin)
			}
			if got := user.IsEditor(); got != tt.isEditor {
				t.Errorf("IsEditor() = %v, want %v", got, tt.isEditor)
			}
			if got := user.IsUser(); got != tt.isUser {
				t.Errorf("IsUser() = %v, want %v", got, tt.isUser)
			}
		})
	}
}

func TestUserToResponse(t *testing.T) {
	user := &User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      RoleUser,
	}
	resp := user.ToResponse()

	if resp.ID != user.ID {
		t.Errorf("ToResponse() ID = %v, want %v", resp.ID, user.ID)
	}
	if resp.Email != user.Email {
		t.Errorf("ToResponse() Email = %v, want %v", resp.Email, user.Email)
	}
	if resp.FirstName != user.FirstName {
		t.Errorf("ToResponse() FirstName = %v, want %v", resp.FirstName, user.FirstName)
	}
	if resp.LastName != user.LastName {
		t.Errorf("ToResponse() LastName = %v, want %v", resp.LastName, user.LastName)
	}
	if resp.Role != user.Role {
		t.Errorf("ToResponse() Role = %v, want %v", resp.Role, user.Role)
	}
	if resp.FullName != user.GetFullName() {
		t.Errorf("ToResponse() FullName = %v, want %v", resp.FullName, user.GetFullName())
	}
}
