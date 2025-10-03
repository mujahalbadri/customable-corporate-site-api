package versions

import (
	"gorm.io/gorm"
)


// Migration version: 002_add_user_indexes
func Migration002AddUserIndexes() MigrationStep {
	return MigrationStep{
		Version:     "002_add_user_indexes",
		Description: "Add indexes to users table",
		Up: func(tx *gorm.DB) error {
			// Add indexes to the users table
			if err := tx.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
				return err
			}

			// Add index on is_active column
			if err := tx.Exec("CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)").Error; err != nil {
				return err
			}

			// Add composite index for active users by role
			if err := tx.Exec("CREATE INDEX IF NOT EXISTS idx_users_is_active_role ON users(is_active, role)").Error; err != nil {
				return err
			}

			return nil
		},
		Down: func(tx *gorm.DB) error {
			// Drop indexes from the users table
			indexes := []string{
				"idx_users_email",
				"idx_users_role",
				"idx_users_is_active",
				"idx_users_is_active_role",
			}

			for _, index := range indexes {
				if err := tx.Exec("DROP INDEX IF EXISTS %s", index).Error; err != nil {
					return err
				}
			}

			return nil
		},
	}
}
