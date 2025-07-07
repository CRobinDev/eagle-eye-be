package gemini

import (
	"context"
	"log"

	"github.com/CRobinDev/karsa/config/env"
	"google.golang.org/genai"
)

func NewGemini() *genai.Client {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  env.GetEnv().GeminiApiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		log.Fatal("failed to create Gemini client:", err)
	}

	return client
}
