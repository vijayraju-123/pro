package config

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB is the global database connection instance
var DB *sql.DB
var once sync.Once // For singleton pattern to ensure DB is initialized only once

// InitDB initializes the PostgreSQL database connection
func InitDB() {
	once.Do(func() {
		var err error
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			ConfigInstance.DBUser, ConfigInstance.DBPassword, ConfigInstance.DBName,
			ConfigInstance.DBHost, ConfigInstance.DBPort)
		log.Printf("Attempting to connect with: %s", connStr) // Debug log

		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}

		// Set connection pool settings for better performance
		DB.SetMaxOpenConns(10)  // Maximum number of open connections
		DB.SetMaxIdleConns(5)   // Maximum number of idle connections
		DB.SetConnMaxLifetime(0) // Connection can live forever (0 = no limit)

		// Test the connection
		err = DB.Ping()
		if err != nil {
			log.Fatalf("Error pinging the database: %v", err)
		}

		log.Println("Database connected successfully!")
	})
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
		log.Println("Database connection closed.")
	}
}

// GetDB returns the database connection (for use in other packages)
func GetDB() *sql.DB {
	if DB == nil {
		InitDB() // Initialize if not already done
	}
	return DB
}