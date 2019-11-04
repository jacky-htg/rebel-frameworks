package main

import (
	"essentials/libraries/config"
	"essentials/libraries/database"
	"essentials/schema"
	"flag"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if _, ok := os.LookupEnv("APP_ENV"); !ok {
		config.Setup(".env")
	}

	// Start Database
	db, err := database.Open()
	if err != nil {
		log.Fatalf("error: connecting to db: %s", err)
	}
	defer db.Close()

	// Handle cli command
	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Println("error applying migrations", err)
			os.Exit(1)
		}
		log.Println("Migrations complete")
		return

	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Println("error seeding database", err)
			os.Exit(1)
		}
		log.Println("Seed data complete")
		return
	}
}
