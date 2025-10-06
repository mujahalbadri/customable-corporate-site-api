package postgres

import (
	"customable-corporate-site-api/internal/models"
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the User model
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Failed to auto-migrate test database: %v", err)
	}

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      models.RoleUser,
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Errorf("Expected user ID to be set, got 0")
	}

	if user.Password == "password123" {
		t.Errorf("Expected password to be hashed, but it was not")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      models.RoleUser,
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve user by ID: %v", err)
	}

	if retrievedUser == nil {
		t.Errorf("Expected user to be found, got nil")
		return
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, retrievedUser.ID)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      models.RoleUser,
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrievedUser, err := repo.GetByEmail(user.Email)
	if err != nil {
		t.Fatalf("Failed to retrieve user by email: %v", err)
	}

	if retrievedUser == nil {
		t.Errorf("Expected user to be found, got nil")
		return
	}

	if retrievedUser.Email != user.Email {
		t.Errorf("Expected user email %q, got %q", user.Email, retrievedUser.Email)
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      models.RoleUser,
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user.FirstName = "Jane"
	if err := repo.Update(user); err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	updatedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve user by ID: %v", err)
	}
	if updatedUser.FirstName != "Jane" {
		t.Errorf("Expected first name to be updated to 'Jane', got %q", updatedUser.FirstName)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      models.RoleUser,
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if err := repo.Delete(user.ID); err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// After delete, GORM should return ErrRecordNotFound when fetching
	deletedUser, err := repo.GetByID(user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// expected: user deleted
		} else {
			t.Fatalf("Failed to retrieve user by ID: %v", err)
		}
	} else {
		if deletedUser != nil {
			t.Errorf("Expected user to be deleted, got %v", deletedUser)
		}
	}
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create multiple users
	users := []models.User{
		{Email: "user1@example.com", Password: "password123", FirstName: "User", LastName: "One", Role: models.RoleUser},
		{Email: "user2@example.com", Password: "password123", FirstName: "User", LastName: "Two", Role: models.RoleUser},
		{Email: "user3@example.com", Password: "password123", FirstName: "User", LastName: "Three", Role: models.RoleUser},
		{Email: "user4@example.com", Password: "password123", FirstName: "User", LastName: "Four", Role: models.RoleUser},
		{Email: "user5@example.com", Password: "password123", FirstName: "User", LastName: "Five", Role: models.RoleUser},
	}

	for i := range users {
		if err := repo.Create(&users[i]); err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	// Test Pagination (offset 0)
	retrievedUsers, err := repo.List(0, 10) // Offset 0, Limit 10
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if len(retrievedUsers) != 5 {
		t.Errorf("Expected 5 users, got %d", len(retrievedUsers))
	}
}

func TestUserRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create test users
	users := []models.User{
		{Email: "user1@example.com", Password: "password123", FirstName: "User", LastName: "One", Role: models.RoleUser},
		{Email: "user2@example.com", Password: "password123", FirstName: "User", LastName: "Two", Role: models.RoleUser},
		{Email: "user3@example.com", Password: "password123", FirstName: "User", LastName: "Three", Role: models.RoleUser},
	}

	for i := range users {
		if err := repo.Create(&users[i]); err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	count, err := repo.Count()
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	if count != int64(len(users)) {
		t.Errorf("Expected user count %d, got %d", len(users), count)
	}
}

func TestUserRepository_GetActiveUsers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create active and inactive users
	activeUser := &models.User{
		Email:     "active@example.com",
		Password:  "password123",
		FirstName: "Active",
		LastName:  "User",
		Role:      models.RoleUser,
		IsActive:  true,
	}

	inactiveUser := &models.User{
		Email:     "inactive@example.com",
		Password:  "password123",
		FirstName: "Inactive",
		LastName:  "User",
		Role:      models.RoleUser,
		IsActive:  false,
	}

	if err := repo.Create(activeUser); err != nil {
		t.Fatalf("Failed to create active user: %v", err)
	}

	if err := repo.Create(inactiveUser); err != nil {
		t.Fatalf("Failed to create inactive user: %v", err)
	}

	// Ensure the second user is marked inactive in DB (some drivers may apply defaults)
	if err := repo.UpdateUserStatus(inactiveUser.ID, false); err != nil {
		t.Fatalf("Failed to mark user inactive: %v", err)
	}

	// GetActiveUsers(limit, offset)
	activeUsers, err := repo.GetActiveUsers(10, 0)
	if err != nil {
		t.Fatalf("Failed to get active users: %v", err)
	}

	if len(activeUsers) != 1 {
		t.Errorf("Expected 1 active user, got %d", len(activeUsers))
	}

	if activeUsers[0].Email != activeUser.Email {
		t.Errorf("Expected active user email %q, got %q", activeUser.Email, activeUsers[0].Email)
	}
}

func TestUserRepository_GetUserByRole(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create users with different roles
	adminUser := &models.User{
		Email:     "admin@example.com",
		Password:  "password123",
		FirstName: "Admin",
		LastName:  "User",
		Role:      models.RoleAdmin,
	}

	editorUser := &models.User{
		Email:     "editor@example.com",
		Password:  "password123",
		FirstName: "Editor",
		LastName:  "User",
		Role:      models.RoleEditor,
	}

	viewerUser := &models.User{
		Email:     "viewer@example.com",
		Password:  "password123",
		FirstName: "Viewer",
		LastName:  "User",
		Role:      models.RoleUser,
	}

	if err := repo.Create(adminUser); err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	if err := repo.Create(editorUser); err != nil {
		t.Fatalf("Failed to create editor user: %v", err)
	}

	if err := repo.Create(viewerUser); err != nil {
		t.Fatalf("Failed to create viewer user: %v", err)
	}

	admins, err := repo.GetUsersByRole(models.RoleAdmin, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get admin users: %v", err)
	}

	if len(admins) != 1 {
		t.Errorf("Expected 1 admin user, got %d", len(admins))
	}

	if admins[0].Role != models.RoleAdmin {
		t.Errorf("Expected admin user role %q, got %q", models.RoleAdmin, admins[0].Role)
	}

	editors, err := repo.GetUsersByRole(models.RoleEditor, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get editor users: %v", err)
	}

	if len(editors) != 1 {
		t.Errorf("Expected 1 editor user, got %d", len(editors))
	}

	if editors[0].Role != models.RoleEditor {
		t.Errorf("Expected editor user role %q, got %q", models.RoleEditor, editors[0].Role)
	}

	viewers, err := repo.GetUsersByRole(models.RoleUser, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get viewer users: %v", err)
	}

	if len(viewers) != 1 {
		t.Errorf("Expected 1 viewer user, got %d", len(viewers))
	}

	if viewers[0].Role != models.RoleUser {
		t.Errorf("Expected viewer user role %q, got %q", models.RoleUser, viewers[0].Role)
	}
}

func TestUserRepository_SearchUsers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create test users
	users := []models.User{
		{Email: "user1@example.com", Password: "password123", FirstName: "User", LastName: "One", Role: models.RoleUser},
		{Email: "user2@example.com", Password: "password123", FirstName: "User", LastName: "Two", Role: models.RoleUser},
		{Email: "user3@example.com", Password: "password123", FirstName: "User", LastName: "Three", Role: models.RoleUser},
	}

	for _, user := range users {
		if err := repo.Create(&user); err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	tests := []struct {
		name      string
		query     string
		wantCount int
	}{
		{"SearchByFirstName", "User", 3},
		{"SearchByLastName", "Two", 1},
		{"SearchByEmail", "user1@example.com", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.SearchUsers(tt.query, 10, 0)
			if err != nil {
				t.Fatalf("SearchUsers() error = %v", err)
			}
			if len(got) != tt.wantCount {
				t.Errorf("SearchUsers() returned %d users, got %d", tt.wantCount, len(got))
			}
		})
	}
}
