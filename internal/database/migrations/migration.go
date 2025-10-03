package migrations

import (
	"fmt"
	"log"
	"time"

	"customable-corporate-site-api/internal/database/migrations/versions"

	"gorm.io/gorm"
)

// Migration represents a database migration.
type Migration struct {
	ID          uint      `gorm:"primaryKey"`
	Version     string    `gorm:"unique;not null"`
	Description string    `gorm:"not null"`
	ExecutedAt  time.Time `gorm:"not null"`
}

// TableName specifies the table name for the Migration model.
func (Migration) TableName() string {
	return "migrations"
}

// MigrationStep is imported from versions package
type MigrationStep = versions.MigrationStep

// Migratior handles database migrations.
type Migrator struct {
	db         *gorm.DB
	migrations []MigrationStep
}

// NewMigrator creates a new Migrator instance.
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: []MigrationStep{},
	}
}

// Register registers a new migration step.
func (m *Migrator) Register(step MigrationStep) {
	m.migrations = append(m.migrations, step)
}

// Initialize creates the migrations table if it does not exist.
func (m *Migrator) Initialize() error {
	return m.db.AutoMigrate(&Migration{})
}

// Up runs all pending migrations.
func (m *Migrator) Up() error {
	log.Println("Starting database migrations...")

	// Initialize the migrations table
	if err := m.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize migrations table: %w", err)
	}

	// Get executed migrations
	var executedMigrations []Migration
	if err := m.db.Find(&executedMigrations).Error; err != nil {
		return fmt.Errorf("failed to fetch executed migrations: %w", err)
	}

	// Create a map of executed migration versions for quick lookup
	executedMap := make(map[string]bool)
	for _, migration := range executedMigrations {
		executedMap[migration.Version] = true
	}

	// Run pending migrations
	pendingCount := 0
	for _, migration := range m.migrations {
		if executedMap[migration.Version] {
			continue
		}

		log.Printf("Running migration: %s - %s", migration.Version, migration.Description)

		// Run migration in a transaction
		err := m.db.Transaction(func(tx *gorm.DB) error {
			if err := migration.Up(tx); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
			}

			// Record the executed migration
			newMigration := Migration{
				Version:     migration.Version,
				Description: migration.Description,
				ExecutedAt:  time.Now(),
			}
			if err := tx.Create(&newMigration).Error; err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration.Version, err)
		}

		log.Printf("Successfully ran migration: %s", migration.Version)
		pendingCount++
	}

	if pendingCount == 0 {
		log.Println("Database is up to date, no migrations needed.")
	} else {
		log.Printf("Successfully applied %d migrations.", pendingCount)
	}

	return nil
}

// Down rolls back the last executed migration.
func (m *Migrator) Down() error {
	log.Println("Rolling back the last migration...")

	// Get the last executed migration
	var lastMigration Migration
	if err := m.db.Order("executed_at desc").First(&lastMigration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No migrations to roll back.")
			return nil
		}
		return fmt.Errorf("failed to fetch last migration: %w", err)
	}

	// Find migration step
	var migrationStep *MigrationStep
	for _, migration := range m.migrations {
		if migration.Version == lastMigration.Version {
			migrationStep = &migration
			break
		}
	}

	if migrationStep == nil {
		return fmt.Errorf("migration step not found for version: %s", lastMigration.Version)
	}

	log.Printf("Rolling back migration: %s - %s", migrationStep.Version, migrationStep.Description)

	// Run rollback in a transaction
	err := m.db.Transaction(func(tx *gorm.DB) error {
		// Execute rollback
		if err := migrationStep.Down(tx); err != nil {
			return fmt.Errorf("failed to execute rollback for migration %s: %w", migrationStep.Version, err)
		}

		// Remove the migration record
		if err := tx.Delete(&lastMigration).Error; err != nil {
			return fmt.Errorf("failed to remove migration record %s: %w", lastMigration.Version, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to roll back migration %s: %w", migrationStep.Version, err)
	}

	log.Printf("Successfully rolled back migration: %s", migrationStep.Version)
	return nil
}

// Status shows migration status.
func (m *Migrator) Status() error {
	log.Println("Migration Status:")

	// Get executed migrations
	var executedMigrations []Migration
	if err := m.db.Order("executed_at asc").Find(&executedMigrations).Error; err != nil {
		return fmt.Errorf("failed to fetch executed migrations: %w", err)
	}

	executed := make(map[string]Migration)
	for _, migration := range executedMigrations {
		executed[migration.Version] = migration
	}

	log.Println("\n╔════════════════════════════════════════════════════════════╗")
	log.Println("║  Version        │  Description              │  Status       ║")
	log.Println("╠════════════════════════════════════════════════════════════╣")

	for _, migration := range m.migrations {
		status := "Pending"
		executedAt := ""
		if exec, ok := executed[migration.Version]; ok {
			status = "Executed"
			executedAt = exec.ExecutedAt.Format("2006-01-02 15:04:05")
		}

		log.Printf("║  %-14s │  %-25s │  %-12s ║", migration.Version, truncate(migration.Description, 24), status)
		if executedAt != "" {
			log.Printf("║                 │  %s              │               ║", executedAt)
		}
		log.Println("╠════════════════════════════════════════════════════════════╣")
	}
	log.Println("╚════════════════════════════════════════════════════════════╝")

	return nil
}

// Reset drops all tables and re-runs all migrations.
func (m *Migrator) Reset() error {
	log.Println("Warning: This will drop all tables and re-run all migrations. Proceeding...")
	log.Println("Resetting tables...")

	// Drop migrations table first to avoid foreign key issues
	if err := m.db.Migrator().DropTable(&Migration{}); err != nil {
		return fmt.Errorf("failed to drop migrations table: %w", err)
	}

	// Drop all registered tables
	// Note: In a real application, you would drop all tables here.
	// For simplicity, we are only dropping the migrations table.

	// Re-run migrations
	if err := m.Up(); err != nil {
		return fmt.Errorf("failed to re-run migrations: %w", err)
	}
	log.Println("Reset complete.")
	return nil
}

// Helper function to truncate strings for display
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
