package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jonkermoo/rag-textbook/backend/internal/models"
	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

// Create a new database connection
func NewDB() (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Finds the most similar chunks to a query embedding
func (db *DB) SearchSimilarChunks(textbookID int, queryEmbedding []float32, topK int) ([]models.Chunk, error) {
	// Convert embedding to pgvector format
	embeddingStr := fmt.Sprintf("[%v]", arrayToString(queryEmbedding))

	query := `
		SELECT id, textbook_id, content, page_number, chunk_index, created_at,
		       embedding <=> $1::vector AS distance
		FROM chunks
		WHERE textbook_id = $2
		ORDER BY embedding <=> $1::vector
		LIMIT $3
	`

	rows, err := db.conn.Query(query, embeddingStr, textbookID, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to query similar chunks: %w", err)
	}
	defer rows.Close()

	var chunks []models.Chunk
	for rows.Next() {
		var chunk models.Chunk
		var distance float64

		err := rows.Scan(
			&chunk.ID,
			&chunk.TextbookID,
			&chunk.Content,
			&chunk.PageNumber,
			&chunk.ChunkIndex,
			&chunk.CreatedAt,
			&distance,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chunk: %w", err)
		}

		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

// Retrieve a textbook by ID
func (db *DB) GetTextbook(id int) (*models.Textbook, error) {
	var textbook models.Textbook

	query := `SELECT id, user_id, title, s3_key, uploaded_at, processed FROM textbooks WHERE id = $1`
	err := db.conn.QueryRow(query, id).Scan(
		&textbook.ID,
		&textbook.UserID,
		&textbook.Title,
		&textbook.S3Key,
		&textbook.UploadedAt,
		&textbook.Processed,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("textbook not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get textbook: %w", err)
	}

	return &textbook, nil
}

// Create a new textbook record
func (db *DB) CreateTextbook(userID int, title, s3Key string) (*models.Textbook, error) {
	var textbook models.Textbook

	query := `
		INSERT INTO textbooks (user_id, title, s3_key, processed)
		VALUES ($1, $2, $3, false)
		RETURNING id, user_id, title, s3_key, uploaded_at, processed
	`

	err := db.conn.QueryRow(query, userID, title, s3Key).Scan(
		&textbook.ID,
		&textbook.UserID,
		&textbook.Title,
		&textbook.S3Key,
		&textbook.UploadedAt,
		&textbook.Processed,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create textbook: %w", err)
	}

	return &textbook, nil
}

// Create a new user with hashed password and verification token
func (db *DB) CreateUser(email, passwordHash, verificationToken string) (*models.User, error) {
	var user models.User

	query := `
		INSERT INTO users (email, password_hash, verification_token, verified)
		VALUES ($1, $2, $3, false)
		RETURNING id, email, verified, created_at
	`

	err := db.conn.QueryRow(query, email, passwordHash, verificationToken).Scan(
		&user.ID,
		&user.Email,
		&user.Verified,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// Retrieve a user by their email address
func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, email, password_hash, verified, created_at, last_login
		FROM users
		WHERE email = $1
	`

	err := db.conn.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Verified,
		&user.CreatedAt,
		&user.LastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// Mark a user's email as verified using their verification token
func (db *DB) VerifyUser(token string) error {
	query := `
		UPDATE users
		SET verified = true, verification_token = NULL
		WHERE verification_token = $1 AND verified = false
	`

	result, err := db.conn.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invalid or already used verification token")
	}

	return nil
}

// Update the user's last login timestamp
func (db *DB) UpdateLastLogin(userID int) error {
	query := `UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1`

	_, err := db.conn.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// Close the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Helper function to convert float slice to string for pgvector
func arrayToString(arr []float32) string {
	result := ""
	for i, v := range arr {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%f", v)
	}
	return result
}
