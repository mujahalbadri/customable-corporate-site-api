package versions

import (
	"customable-corporate-site-api/internal/models"
	"time"

	"gorm.io/gorm"
)

// Migration version: 003_seed_admin_user
func Migration003SeedAdminUser() MigrationStep {
	return MigrationStep{
		Version:     "003_seed_admin_user",
		Description: "Seed initial admin user",
		Up: func(tx *gorm.DB) error {
			// Check if the admin user already exists
			var count int64
			if err := tx.Model(&models.User{}).Where("email = ?", models.RoleAdmin).Count(&count).Error; err != nil {
				return err
			}

			// Skip if the admin user already exists
			if count > 0 {
				return nil
			}

			// Create default admin user
			adminUser := models.User{
				Email:     "admin@company.com",
				Password:  "Admin123#",
				FirstName: "John",
				LastName:  "Doe",
				Role:      models.RoleAdmin,
				IsActive:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			return tx.Create(&adminUser).Error
		},
		Down: func(tx *gorm.DB) error {
			// Remove the admin user
			return tx.Where("email = ?", "admin@company.com").Delete(&models.User{}).Error
		},
	}
}
