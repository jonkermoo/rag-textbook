package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/jonkermoo/rag-textbook/backend/internal/database"
	"github.com/jonkermoo/rag-textbook/backend/internal/handlers"
	"github.com/jonkermoo/rag-textbook/backend/internal/middleware"
	"github.com/jonkermoo/rag-textbook/backend/internal/services"
)

func main() {
	// Load environment variables
	// Try loading from multiple possible locations (for local dev vs Docker)
	err := godotenv.Load("../../.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Warning: Could not load .env file, using environment variables")
		}
	}

	// Initialize database connection
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Initialize services
	embeddingService := services.NewEmbeddingService()
	log.Println("Embedding service initialized")

	ragService := services.NewRAGService(db, embeddingService)
	log.Println("RAG service initialized")

	authService := services.NewAuthService(db)
	log.Println("Auth service initialized")

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(authService)

	// CORS middleware wrapper
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow multiple origins for development and production
			origin := r.Header.Get("Origin")
			allowedOrigins := []string{
				"http://localhost:5173",    // Local development
				"http://localhost:3000",    // Alternative local port
				"https://www.lexra.online", // Production frontend (www)
				"https://lexra.online",     // Production frontend (no www)
				"https://lexra-1539-jwbn-johnny-zhus-projects-0f83c3d9.vercel.app", // Vercel deployment
			}

			// Check if origin is in allowed list
			isAllowed := false
			for _, allowed := range allowedOrigins {
				if origin == allowed || allowed == "*" {
					isAllowed = true
					break
				}
			}

			// Also allow any Vercel preview URL or lexra.online subdomain
			if strings.HasSuffix(origin, ".vercel.app") || strings.Contains(origin, "lexra.online") {
				isAllowed = true
			}

			if isAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Initialize handlers
	textbookHandler := handlers.NewTextbookHandler(db)
	queryHandler := handlers.NewQueryHandler(ragService)
	authHandler := handlers.NewAuthHandler(authService)
	uploadHandler := handlers.NewUploadHandler(db)

	// Textbook management routes (protected)
	http.Handle("/api/textbooks", corsMiddleware(authMiddleware(http.HandlerFunc(textbookHandler.HandleListTextbooks))))
	http.HandleFunc("/api/textbooks/", func(w http.ResponseWriter, r *http.Request) {
		// Apply CORS and auth middleware manually for paths with IDs
		corsMiddleware(authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Route to appropriate handler based on path
			if strings.HasSuffix(r.URL.Path, "/status") {
				textbookHandler.HandleGetTextbookStatus(w, r)
			} else if r.Method == http.MethodDelete {
				textbookHandler.HandleDeleteTextbook(w, r)
			} else if r.Method == http.MethodGet {
				textbookHandler.HandleGetTextbook(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}))).ServeHTTP(w, r)
	})

	// Set up HTTP routes
	// Protected routes (require authentication)
	http.Handle("/api/query", corsMiddleware(authMiddleware(http.HandlerFunc(queryHandler.HandleQuery))))
	http.Handle("/api/upload", corsMiddleware(authMiddleware(http.HandlerFunc(uploadHandler.HandleUpload))))

	// Public routes (no authentication needed)
	http.Handle("/api/auth/register", corsMiddleware(http.HandlerFunc(authHandler.HandleRegister)))
	http.Handle("/api/auth/login", corsMiddleware(http.HandlerFunc(authHandler.HandleLogin)))
	http.Handle("/api/auth/verify", corsMiddleware(http.HandlerFunc(authHandler.HandleVerify)))
	http.Handle("/api/health", corsMiddleware(http.HandlerFunc(handlers.HandleHealth)))

	// Enable CORS for frontend
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 404 for unknown routes
		http.NotFound(w, r)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("\nServer starting on http://localhost:%s", port)
	log.Println("\nAvailable endpoints:")
	log.Println("  POST /api/upload - Upload a textbook PDF")
	log.Println("  GET    /api/textbooks       - List user's textbooks")
	log.Println("  GET    /api/textbooks/:id   - Get textbook details")
	log.Println("  DELETE /api/textbooks/:id   - Delete a textbook")
	log.Println("  GET    /api/textbooks/:id/status - Get processing status")
	log.Println("  POST /api/query  - Submit a question")
	log.Println("  GET  /api/health - Health check")
	log.Println("\nPress Ctrl+C to stop")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
