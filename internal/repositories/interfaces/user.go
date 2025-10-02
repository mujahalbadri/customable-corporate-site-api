package interfaces

import "customable-corporate-site-api/internal/models"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Basic CRUD operations
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error

	// Query operations
	List(offset, limit int) ([]models.User, error)
	Count() (int64, error)

	// Advanced queries
	GetActiveUsers(limit, offset int) ([]models.User, error)
	GetUsersByRole(role string, limit, offset int) ([]models.User, error)
	SearchUsers(query string, limit, offset int) ([]models.User, error)

	// Bulk operations
	UpdateUserStatus(id uint, isActive bool) error
	UpdateUserRole(id uint, role string) error
}
