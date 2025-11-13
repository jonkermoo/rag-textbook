package services

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type EmbeddingService struct {
	client *openai.Client
}

// Create a new embedding service
func NewEmbeddingService() *EmbeddingService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	return &EmbeddingService{
		client: client,
	}
}

// Converts text to a vector embedding
func (s *EmbeddingService) GenerateEmbedding(text string) ([]float32, error) {
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.SmallEmbedding3,
	}

	resp, err := s.client.CreateEmbeddings(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return resp.Data[0].Embedding, nil
}
