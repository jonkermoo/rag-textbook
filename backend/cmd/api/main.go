package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Dot environment
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Build database connection string
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL")

	// Insert a test user
	email := "alice@example.com"
	var userID int
	err = db.QueryRow(
		"INSERT INTO users (email) VALUES ($1) ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email RETURNING id",
		email,
	).Scan(&userID)

	if err != nil {
		log.Fatal("Failed to insert user:", err)
	}

	fmt.Printf("Inserted/Updated user with ID: %d\n", userID)

	// Query all users
	rows, err := db.Query("SELECT id, email, created_at FROM users")
	if err != nil {
		log.Fatal("Failed to query users:", err)
	}
	defer rows.Close()

	fmt.Println("\nAll users in database:")
	fmt.Println("----------------------------")

	for rows.Next() {
		var id int
		var email string
		var createdAt string

		err := rows.Scan(&id, &email, &createdAt)
		if err != nil {
			log.Fatal("Failed to scan row:", err)
		}

		fmt.Printf("ID: %d | Email: %s | Created: %s\n", id, email, createdAt)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error iterating rows:", err)
	}

	fmt.Println("\n Database test completed successfully!")
}
