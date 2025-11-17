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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jonkermoo/rag-textbook/backend/internal/database"
	"github.com/jonkermoo/rag-textbook/backend/internal/middleware"
	"github.com/jonkermoo/rag-textbook/backend/internal/models"
)

type UploadHandler struct {
	db        *database.DB
	s3Client  *s3.S3
	s3Bucket  string
}

func NewUploadHandler(db *database.DB) *UploadHandler {
	// Initialize AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))

	return &UploadHandler{
		db:       db,
		s3Client: s3.New(sess),
		s3Bucket: os.Getenv("S3_BUCKET_NAME"),
	}
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

	// Generate unique S3 key
	s3Key := fmt.Sprintf("textbooks/%d/%s", userID, header.Filename)

	// Upload file to S3
	_, err = h.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(h.s3Bucket),
		Key:         aws.String(s3Key),
		Body:        file,
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		log.Printf("Failed to upload to S3: %v", err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	log.Printf("File uploaded to S3: %s", s3Key)

	// Get title from form or use filename
	title := r.FormValue("title")
	if title == "" {
		title = strings.TrimSuffix(header.Filename, ".pdf")
	}

	// Create textbook record in database with S3 key
	textbook, err := h.db.CreateTextbook(userID, title, s3Key)
	if err != nil {
		log.Printf("Failed to create textbook record: %v", err)
		http.Error(w, "Failed to create textbook record", http.StatusInternalServerError)
		return
	}

	log.Printf("File uploaded successfully: %s (textbook_id=%d)", s3Key, textbook.ID)
	// Trigger background processing
	h.triggerProcessing(textbook.ID, s3Key)
	log.Printf("Processing triggered for textbook_id=%d", textbook.ID)

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
func (h *UploadHandler) triggerProcessing(textbookID int, s3Key string) {
	go func() {
		log.Printf("Starting background processing for textbook %d", textbookID)

		// Download PDF from S3 to temporary file
		tmpFile := fmt.Sprintf("/tmp/textbook_%d.pdf", textbookID)

		result, err := h.s3Client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(h.s3Bucket),
			Key:    aws.String(s3Key),
		})
		if err != nil {
			log.Printf("Error downloading from S3: %v", err)
			return
		}
		defer result.Body.Close()

		// Create temporary file
		outFile, err := os.Create(tmpFile)
		if err != nil {
			log.Printf("Error creating temp file: %v", err)
			return
		}
		defer outFile.Close()
		defer os.Remove(tmpFile) // Clean up after processing

		// Copy S3 object to file
		_, err = io.Copy(outFile, result.Body)
		if err != nil {
			log.Printf("Error saving temp file: %v", err)
			return
		}

		pythonScript := "../../ingestion/src/process_existing.py"

		// Run the Python script with local temp file
		cmd := exec.Command("python", pythonScript,
			fmt.Sprintf("%d", textbookID),
			tmpFile)

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
