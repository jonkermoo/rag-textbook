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

// helper: central place to decide which origins are allowed
func isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	// 1) Local dev: any localhost port
	if strings.HasPrefix(origin, "http://localhost") {
		return true
	}

	// 2) Main production domains
	if origin == "https://lexra.online" || origin == "https://www.lexra.online" {
		return true
	}

	// 3) Any Vercel deployment/preview
	if strings.HasSuffix(origin, ".vercel.app") {
		return true
	}

	// 4) Any lexra.online subdomain (e.g. https://api.lexra.online)
	if strings.HasSuffix(origin, ".lexra.online") {
		return true
	}

	return false
}

func main() {
	// Load environment variables
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
			origin := r.Header.Get("Origin")

			if isOriginAllowed(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Preflight
			if r.Method == http.MethodOptions {
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
		corsMiddleware(authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	// Protected routes
	http.Handle("/api/query", corsMiddleware(authMiddleware(http.HandlerFunc(queryHandler.HandleQuery))))
	http.Handle("/api/upload", corsMiddleware(authMiddleware(http.HandlerFunc(uploadHandler.HandleUpload))))

	// Public routes
	http.Handle("/api/auth/register", corsMiddleware(http.HandlerFunc(authHandler.HandleRegister)))
	http.Handle("/api/auth/login", corsMiddleware(http.HandlerFunc(authHandler.HandleLogin)))
	http.Handle("/api/auth/verify", corsMiddleware(http.HandlerFunc(authHandler.HandleVerify)))
	http.Handle("/api/health", corsMiddleware(http.HandlerFunc(handlers.HandleHealth)))

	// Fallback for unknown routes – no special CORS needed here
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("\nServer starting on http://localhost:%s", port)
	log.Println("\nAvailable endpoints:")
	log.Println("  POST   /api/upload                 - Upload a textbook PDF")
	log.Println("  GET    /api/textbooks              - List user's textbooks")
	log.Println("  GET    /api/textbooks/:id          - Get textbook details")
	log.Println("  DELETE /api/textbooks/:id          - Delete a textbook")
	log.Println("  GET    /api/textbooks/:id/status   - Get processing status")
	log.Println("  POST   /api/query                  - Submit a question")
	log.Println("  GET    /api/health                 - Health check")
	log.Println("\nPress Ctrl+C to stop")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
