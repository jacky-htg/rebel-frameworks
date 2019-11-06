package main

import (
	"essentials/libraries/config"
	"essentials/libraries/database"
	"essentials/schema"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	if _, ok := os.LookupEnv("APP_ENV"); !ok {
		config.Setup(".env")
	}

	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run() error {

	// Start Database
	db, err := database.Open()
	if err != nil {
		return fmt.Errorf("connecting to db: %v", err)
	}
	defer db.Close()

	// Handle cli command
	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			return fmt.Errorf("applying migrations: %v", err)
		}
		log.Println("Migrations complete")

	case "seed":
		if err := schema.Seed(db); err != nil {
			return fmt.Errorf("seeding database: %v", err)
		}
		log.Println("Seed data complete")
	}

	return nil
}
