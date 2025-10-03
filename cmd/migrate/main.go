package main

import (
	"customable-corporate-site-api/internal/config"
	"customable-corporate-site-api/internal/database"
	"customable-corporate-site-api/internal/database/migrations"

	"flag"
	"fmt"
	"log"	
	"os"
	"strings"
)

func main() {
	// Define command-line flags and initialize the application
	upCmd := flag.Bool("up", false, "Run all pending migrations")
	downCmd := flag.Bool("down", false, "Roll back the last migration")
	statusCmd := flag.Bool("status", false, "Show migration status")
	resetCmd := flag.Bool("reset", false, "Reset the database (WARNING: drops all data)")
	createCmd := flag.String("create", "", "Create a new migration file ")

	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Connect to the database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Register migrations
	migrator := migrations.RegisterMigrations(db)

	// Execute the appropriate command
	switch {
	case *upCmd:
		if err := migrator.Up(); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}

	case *downCmd:
		if err := migrator.Down(); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}

	case *statusCmd:
		if err := migrator.Status(); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

	case *resetCmd:
		log.Println("WARNING: This will drop all data in the database. Are you sure? (y/N)")

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" {
			log.Println("Reset cancelled.")
			os.Exit(0)
		}

		if err := migrator.Reset(); err != nil {
			log.Fatalf("Database reset failed: %v", err)
		}

	case *createCmd != "":
		log.Printf("Creating migration: %s\n", *createCmd)
		// Implement migration file creation
		log.Println("Migration file creation not implemented yet")

	default:
		flag.Usage()
		os.Exit(1)
	}

	log.Println("Operation completed successfully.")
}
