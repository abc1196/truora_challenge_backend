package utils

import (
	"log"

	// Import GORM-related packages.
	"github.com/jinzhu/gorm"
	// Import GORM-postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"../data"
)

// DB - Object related to API Database
var DB *gorm.DB

// Connect - Initializes the CockroachDB
func Connect() {
	// Connect to the "db_troura" database as the "abc11" user.
	const addr = "postgresql://abc11@localhost:26257/db_troura?sslmode=disable"
	var err error
	DB, err = gorm.Open("postgres", addr)
	if err != nil {
		log.Fatal(err)
	}

	// Automatically create the "domains" table based on the Domain model.
	DB.AutoMigrate(&data.Domain{})

	// Automatically create the "servers" table based on the Server model.
	DB.AutoMigrate(&data.Server{})
}
