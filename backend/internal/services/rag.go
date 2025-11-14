package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jonkermoo/rag-textbook/backend/internal/database"
	"github.com/jonkermoo/rag-textbook/backend/internal/models"
	"github.com/sashabaranov/go-openai"
)

type RAGService struct {
	db               *database.DB
	embeddingService *EmbeddingService
	openaiClient     *openai.Client
}

// Create a new RAG service
func NewRAGService(db *database.DB, embeddingService *EmbeddingService) *RAGService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	return &RAGService{
		db:               db,
		embeddingService: embeddingService,
		openaiClient:     client,
	}
}

// perform the complete RAG pipeline
func (s *RAGService) Query(req models.QueryRequest, userID int) (*models.QueryResponse, error) {
	startTime := time.Now()

	// Validate textbook exists
	textbook, err := s.db.GetTextbook(req.TextbookID)
	if err != nil {
		return nil, fmt.Errorf("textbook not found: %w", err)
	}

	// Check if user owns this textbook
	if textbook.UserID != userID {
		return nil, fmt.Errorf("permission denied: you don't own this textbook")
	}

	if !textbook.Processed {
		return nil, fmt.Errorf("textbook not yet processed")
	}

	// Set default topK if not provided
	if req.TopK == 0 {
		req.TopK = 5
	}

	// Convert question to embedding
	queryEmbedding, err := s.embeddingService.GenerateEmbedding(req.Question)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Retrieve similar chunks from database
	chunks, err := s.db.SearchSimilarChunks(req.TextbookID, queryEmbedding, req.TopK)
	if err != nil {
		return nil, fmt.Errorf("failed to search chunks: %w", err)
	}

	// Note: We allow processing even if no chunks found - GPT can still provide general guidance
	// or let the user know the topic isn't covered in the textbook

	// Build context from chunk
	contextStr := buildContext(chunks)

	// Generate answer using GPT-4
	answer, err := s.generateAnswer(req.Question, contextStr, textbook.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to generate answer: %w", err)
	}

	// Build response with sources
	sources := make([]models.ChunkSource, len(chunks))
	for i, chunk := range chunks {
		sources[i] = models.ChunkSource{
			PageNumber: chunk.PageNumber,
			Content:    truncateContent(chunk.Content, 200),
			Similarity: 0.0, // calculate if needed
		}
	}

	timeTaken := time.Since(startTime).Milliseconds()

	return &models.QueryResponse{
		Answer:    answer,
		Sources:   sources,
		Question:  req.Question,
		TimeTaken: float64(timeTaken),
	}, nil
}

// Call GPT-4 to generate an answer
func (s *RAGService) generateAnswer(question, contextStr, textbookTitle string) (string, error) {
	systemPrompt := fmt.Sprintf(`You are a knowledgeable tutor with expertise in the subject matter covered in "%s".

Your task is to answer the student's question using the provided textbook context.

Guidelines:
1. Provide clear, direct answers based on the context provided
2. When relevant information is available, explain the concept thoroughly
3. Include page number citations when referencing specific information (e.g., "According to page 42...")
4. If the exact topic isn't covered in the provided context but you can make a reasonable inference from related content, do so confidently
5. If the question is completely outside the scope of the textbook, politely explain that this topic isn't covered in this particular textbook
6. Use clear, student-friendly language
7. Be confident in your explanations - avoid phrases like "the textbook doesn't explicitly say" or "it's not directly stated"`, textbookTitle)

	userPrompt := fmt.Sprintf(`Context from textbook:
---
%s
---

Student question: %s

Please provide a helpful answer based on the context above.`, contextStr, question)

	resp, err := s.openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   500,
		},
	)

	if err != nil {
		return "", fmt.Errorf("openai api error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from openai")
	}

	return resp.Choices[0].Message.Content, nil
}

// Concatenate chunk content for the prompt
func buildContext(chunks []models.Chunk) string {
	var builder strings.Builder

	for i, chunk := range chunks {
		builder.WriteString(fmt.Sprintf("[Page %d]\n%s\n\n", chunk.PageNumber, chunk.Content))
		if i < len(chunks)-1 {
			builder.WriteString("---\n\n")
		}
	}

	return builder.String()
}

// Limit content length for source display
func truncateContent(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}
