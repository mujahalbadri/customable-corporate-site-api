package migrations

import (
	"customable-corporate-site-api/internal/database/migrations/versions"

	"gorm.io/gorm"
)

// RegisterMigrations registers all migration steps.
func RegisterMigrations(db *gorm.DB) *Migrator {
	migrator := NewMigrator(db)

	// Register migration in the order they should be applied
	migrator.Register(versions.Migration001CreateUsersTable())
	migrator.Register(versions.Migration002AddUserIndexes())
	migrator.Register(versions.Migration003SeedAdminUser())

	return migrator
}
