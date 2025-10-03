package versions

import (
	"customable-corporate-site-api/internal/models"

	"gorm.io/gorm"
)

// MigrationStep represents a single migration step
type MigrationStep struct {
	Version     string
	Description string
	Up          func(tx *gorm.DB) error
	Down        func(tx *gorm.DB) error
}

// Migration version: 001_create_users_table
func Migration001CreateUsersTable() MigrationStep {
	return MigrationStep{
		Version:     "001_create_users_table",
		Description: "Create users table",
		Up: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.User{})
		},
		Down: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&models.User{})
		},
	}
}
