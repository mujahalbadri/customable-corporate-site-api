package postgres

import (
	"customable-corporate-site-api/internal/models"
	"customable-corporate-site-api/internal/repositories/interfaces"
	"strings"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID from the database
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email from the database
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user in the database
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user from the database
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List retrieves a list of users from the database with pagination
func (r *userRepository) List(offset, limit int) ([]models.User, error) {
	var users []models.User
	if err := r.db.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Count returns the total number of users in the database
func (r *userRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetActiveUsers retrieves all active users from the database
func (r *userRepository) GetActiveUsers(limit, offset int) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("is_active = ?", true).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUsersByRole retrieves users by their role from the database
func (r *userRepository) GetUsersByRole(role string, limit, offset int) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("role = ?", role).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// SearchUsers searches users by name or email in the database
func (r *userRepository) SearchUsers(query string, limit, offset int) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(email) LIKE ?", "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUserStatus updates the active status of a user
func (r *userRepository) UpdateUserStatus(id uint, isActive bool) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Update("is_active", isActive).Error; err != nil {
		return err
	}
	return nil
}

// UpdateUserRole updates the role of a user
func (r *userRepository) UpdateUserRole(id uint, role string) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Update("role", role).Error; err != nil {
		return err
	}
	return nil
}
