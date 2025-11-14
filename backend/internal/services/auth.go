package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/jonkermoo/rag-textbook/backend/internal/database"
	"github.com/jonkermoo/rag-textbook/backend/internal/models"
)

type AuthService struct {
	db        *database.DB
	jwtSecret []byte
}

func NewAuthService(db *database.DB) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable not set")
	}

	return &AuthService{
		db:        db,
		jwtSecret: []byte(secret),
	}
}

// Create a new user account
func (s *AuthService) Register(req models.RegisterRequest) (*models.User, error) {
	// Validate email
	if !isValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate password strength
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate verification token
	verificationToken, err := generateRandomToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Create user in database
	user, err := s.db.CreateUser(req.Email, string(hashedPassword), verificationToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// TODO: Send verification email (we'll skip for now)
	// For development, auto-verify
	if os.Getenv("AUTO_VERIFY") == "true" {
		user.Verified = true
		s.db.VerifyUser(verificationToken)
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req models.LoginRequest) (string, *models.User, error) {
	// Get user from database
	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Check if verified
	if !user.Verified {
		return "", nil, errors.New("email not verified")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateJWT(user.ID, user.Email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login
	s.db.UpdateLastLogin(user.ID)

	return token, user, nil
}

// Mark a user's email as verified
func (s *AuthService) VerifyEmail(token string) error {
	return s.db.VerifyUser(token)
}

// Validate a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

// Create a new JWT token
func (s *AuthService) generateJWT(userID int, email string) (string, error) {
	expirationHours := 168 // 7 days
	if envHours := os.Getenv("JWT_EXPIRATION_HOURS"); envHours != "" {
		fmt.Sscanf(envHours, "%d", &expirationHours)
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// Helper functions
func isValidEmail(email string) bool {
	// Basic email validation
	return len(email) > 3 && contains(email, "@") && contains(email, ".")
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
