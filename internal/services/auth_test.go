package services

import (
	"customable-corporate-site-api/internal/models"
	"customable-corporate-site-api/internal/repositories/postgres"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestService(t *testing.T) (*AuthService, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	userRepo := postgres.NewUserRepository(db)
	authService := NewAuthService(userRepo, "test_secret-key", 24*time.Hour)

	return authService, db
}

func TestAuthService_Register(t *testing.T) {
	authService, _ := setupTestService(t)

	// Test cases
	tests := []struct {
		name    string
		req     *RegisterRequest
		wantErr bool
	}{
		{
			name: "Valid registration",
			req: &RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: false,
		},
		{
			name: "Duplicate email",
			req: &RegisterRequest{
				Email:     "test@example.com",
				Password:  "newpassword",
				FirstName: "Jane",
				LastName:  "Smith",
			},
			wantErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := authService.Register(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil || resp.User == nil {
					t.Errorf("Register() got nil response or user")
				}
				if resp.User.Email != tt.req.Email {
					t.Errorf("Register() got user email = %v, want %v", resp.User.Email, tt.req.Email)
				}
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	authService, _ := setupTestService(t)

	// Create a test user
	registerReq := &RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	_, err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	// Test cases
	tests := []struct {
		name    string
		req     *LoginRequest
		wantErr bool
	}{
		{
			name: "Valid login",
			req: &LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Invalid password",
			req: &LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			wantErr: true,
		},
		{
			name: "Non-existent email",
			req: &LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := authService.Login(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil || resp.Token == nil {
					t.Errorf("Login() got nil response or token")
				}
				if resp.Token.AccessToken == "" {
					t.Errorf("Login() got empty access token")
				}
				if resp.Token.RefreshToken == "" {
					t.Errorf("Login() got empty refresh token")
				}
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	authService, _ := setupTestService(t)

	// Create a test user and login to get a token
	registerReq := &RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	_, err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	loginReq := &LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	loginResp, err := authService.Login(loginReq)
	if err != nil {
		t.Fatalf("Failed to login test user: %v", err)
	}

	// Test cases
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Valid token",
			token:   loginResp.Token.AccessToken,
			wantErr: false,
		},
		{
			name:    "Invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authService.ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && claims == nil {
				t.Errorf("ValidateToken() got nil claims")
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	authService, _ := setupTestService(t)

	// Create a test user and login to get a token
	registerReq := &RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	_, err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	loginReq := &LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	loginResp, err := authService.Login(loginReq)
	if err != nil {
		t.Fatalf("Failed to login test user: %v", err)
	}

	// Test cases
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Valid refresh token",
			token:   loginResp.Token.RefreshToken,
			wantErr: false,
		},
		{
			name:    "Invalid refresh token",
			token:   "invalid.refresh.token",
			wantErr: true,
		},
		{
			name:    "Empty refresh token",
			token:   "",
			wantErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := authService.RefreshToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp.AccessToken == "" {
					t.Errorf("RefreshToken() got empty access token")
				}
				if resp.RefreshToken == "" {
					t.Errorf("RefreshToken() got empty refresh token")
				}
			}
		})
	}
}

func TestAuthService_GetProfile(t *testing.T) {
	authService, _ := setupTestService(t)

	// Create a test user
	registerReq := &RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	registerResp, err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	// Test cases
	tests := []struct {
		name    string
		userID  uint
		wantErr bool
	}{
		{
			name:    "Valid user ID",
			userID:  registerResp.User.ID,
			wantErr: false,
		},
		{
			name:    "Non-existent user ID",
			userID:  9999,
			wantErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := authService.GetProfile(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if user == nil {
					t.Errorf("GetProfile() got nil user")
				}
				if user != nil && (user.ID != tt.userID) {
					t.Errorf("GetProfile() got user ID = %v, want %v", user.ID, tt.userID)
				}
			}
		})
	}
}

func TestAuthService_UpdateProfile(t *testing.T) {
	authService, _ := setupTestService(t)

	// Create a test user
	registerReq := &RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	_, err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}

	// Test cases
	tests := []struct {
		name    string
		req     *UpdateProfileRequest
		wantErr bool
	}{
		{
			name: "Valid update",
			req: &UpdateProfileRequest{
				FirstName: "Jane",
				LastName:  "Smith",
			},
			wantErr: false,
		},
		{
			name: "Non-existent user ID",
			req: &UpdateProfileRequest{
				FirstName: "No",
				LastName:  "Body",
			},
			wantErr: true,
		},
	}
	// Run tests
	// First case: update the created user (ID should be 1)
	t.Run(tests[0].name, func(t *testing.T) {
		// Get created user
		user, err := authService.GetProfile(1)
		if err != nil {
			t.Fatalf("Failed to get created user: %v", err)
		}
		updated, err := authService.UpdateProfile(user.ID, tests[0].req)
		if (err != nil) != tests[0].wantErr {
			t.Errorf("UpdateProfile() error = %v, wantErr %v", err, tests[0].wantErr)
			return
		}
		if !tests[0].wantErr {
			if updated == nil {
				t.Errorf("UpdateProfile() got nil user")
			}
			if updated != nil && (updated.FirstName != tests[0].req.FirstName || updated.LastName != tests[0].req.LastName) {
				t.Errorf("UpdateProfile() got user name = %v %v, want %v %v", updated.FirstName, updated.LastName, tests[0].req.FirstName, tests[0].req.LastName)
			}
		}
	})

	// Second case: attempt update on non-existent user ID
	t.Run(tests[1].name, func(t *testing.T) {
		_, err := authService.UpdateProfile(9999, tests[1].req)
		if (err != nil) != tests[1].wantErr {
			t.Errorf("UpdateProfile() error = %v, wantErr %v", err, tests[1].wantErr)
		}
	})
}
