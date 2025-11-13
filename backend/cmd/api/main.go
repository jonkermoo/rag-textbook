package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/jonkermoo/rag-textbook/backend/internal/database"
	"github.com/jonkermoo/rag-textbook/backend/internal/handlers"
	"github.com/jonkermoo/rag-textbook/backend/internal/services"
)

func main() {
	// Load environment variables
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Initialize database connection
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	log.Println("✓ Connected to database")

	// Initialize services
	embeddingService := services.NewEmbeddingService()
	log.Println("✓ Embedding service initialized")

	ragService := services.NewRAGService(db, embeddingService)
	log.Println("✓ RAG service initialized")

	// Initialize handlers
	queryHandler := handlers.NewQueryHandler(ragService)

	// Set up HTTP routes
	http.HandleFunc("/api/query", queryHandler.HandleQuery)
	http.HandleFunc("/api/health", handlers.HandleHealth)

	// Enable CORS for frontend
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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
	log.Println("  POST /api/query  - Submit a question")
	log.Println("  GET  /api/health - Health check")
	log.Println("\nPress Ctrl+C to stop")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
