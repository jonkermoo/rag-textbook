package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jonkermoo/rag-textbook/backend/internal/database"
	"github.com/jonkermoo/rag-textbook/backend/internal/middleware"
	"github.com/jonkermoo/rag-textbook/backend/internal/models"
)

type UploadHandler struct {
	db *database.DB
}

func NewUploadHandler(db *database.DB) *UploadHandler {
	return &UploadHandler{db: db}
}

// Handle PDF upload
func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (added by auth middleware)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form (max 50MB)
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
		return
	}

	// Get the file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file is a PDF
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Printf("Failed to create uploads directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", userID, header.Filename)
	filepath := filepath.Join(uploadsDir, filename)

	// Save file to disk
	dst, err := os.Create(filepath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Failed to copy file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Get title from form or use filename
	title := r.FormValue("title")
	if title == "" {
		title = strings.TrimSuffix(header.Filename, ".pdf")
	}

	// Create textbook record in database
	textbook, err := h.db.CreateTextbook(userID, title, filepath)
	if err != nil {
		log.Printf("Failed to create textbook record: %v", err)
		http.Error(w, "Failed to create textbook record", http.StatusInternalServerError)
		return
	}

	log.Printf("File uploaded successfully: %s (textbook_id=%d)", filename, textbook.ID)
	// Trigger background processing
	h.triggerProcessing(textbook.ID, filepath)
	log.Printf("File uploaded successfully: %s (textbook_id=%d)", filename, textbook.ID)

	// Return response
	response := models.UploadResponse{
		TextbookID: textbook.ID,
		Title:      textbook.Title,
		Message:    "File uploaded successfully. Processing will begin shortly.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// triggerProcessing launches the Python ingestion pipeline in the background
func (h *UploadHandler) triggerProcessing(textbookID int, pdfPath string) {
	go func() {
		log.Printf("Starting background processing for textbook %d", textbookID)

		pythonScript := "../../ingestion/src/process_existing.py"

		// Run the Python script
		cmd := exec.Command("python", pythonScript,
			fmt.Sprintf("%d", textbookID),
			pdfPath)

		// Capture output
		output, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("Error processing textbook %d: %v", textbookID, err)
			log.Printf("Python output: %s", string(output))
		} else {
			log.Printf("Successfully processed textbook %d", textbookID)
			log.Printf("Python output: %s", string(output))
		}
	}()
}
