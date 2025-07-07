package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/CRobinDev/karsa/config/env"
	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/sirupsen/logrus"
	"google.golang.org/genai"
)

type IGemini interface {
	DetectDeepFakeImage(ctx context.Context, req dto.GeminiAnalyzeRequest) (dto.GeminiAnalyzeResponse, error)
	DetectDeepFakeAudio(ctx context.Context, req dto.GeminiAnalyzeRequest) (dto.GeminiAnalyzeResponse, error)
}

type gemini struct {
	client *genai.Client
	logger *logrus.Logger
}

func NewGeminiService(client *genai.Client, logger *logrus.Logger) IGemini {
	return &gemini{
		client: client,
		logger: logger,
	}
}

func (g *gemini) DetectDeepFakeImage(ctx context.Context, req dto.GeminiAnalyzeRequest) (dto.GeminiAnalyzeResponse, error) {

	prompt := []*genai.Part{
		{
			Text: `Respond this message as an DeepFake Detector Expert. I will provide you an user image. 
					I want you to predict if the image is DeepFake detected. Please also provide the confidence score. 
						Please respond only in JSON format as shown below :   
						{
							"prediction": "fake",
							"confidence": 0.914
						}
					or 
						{
							"prediction": "real",
							"confidence" : 0.832
						}.
					`,
		},
	}

	return g.processAnalyzation(ctx, req, prompt)
}

func (g *gemini) DetectDeepFakeAudio(ctx context.Context, req dto.GeminiAnalyzeRequest) (dto.GeminiAnalyzeResponse, error) {

	prompt := []*genai.Part{
		{
			Text: `Respond this message as an DeepFake Detector Expert. I will provide you an user audio. 
					I want you to predict if the audio is DeepFake detected. Please also provide the confidence score. 
						Please respond only in JSON format as shown below :   
						{
							"prediction": "fake",
							"confidence": 0.914
						}
					or 
						{
							"prediction": "real",
							"confidence" : 0.832
						}.
					`,
		},
	}

	return g.processAnalyzation(ctx, req, prompt)
}

func (g *gemini) processAnalyzation(ctx context.Context, req dto.GeminiAnalyzeRequest, prompt []*genai.Part) (dto.GeminiAnalyzeResponse, error) {
	prompt = append(prompt, &genai.Part{
		InlineData: &genai.Blob{
			Data:     req.File,
			MIMEType: req.ContentType,
		},
	})

	result, err := g.client.Models.GenerateContent(ctx, env.GetEnv().GeminiApiModel, []*genai.Content{{Parts: prompt}}, nil)
	if err != nil {
		return dto.GeminiAnalyzeResponse{}, fmt.Errorf("failed to generate content: %w", err)
	}

	part := result.Candidates[0].Content.Parts[0]
	
	rawJSON := []byte(part.Text)

	cleanedJSON := cleanJSON(string(rawJSON))

	var resp dto.GeminiAnalyzeResponse

	if err = json.Unmarshal([]byte(cleanedJSON), &resp); err != nil {
		return dto.GeminiAnalyzeResponse{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	log.Println(cleanedJSON)

	return resp, nil
}

func cleanJSON(input string) string {
	cleaned := strings.TrimSpace(input)
	cleaned = strings.Trim(cleaned, "`")
	cleaned = strings.TrimPrefix(cleaned, "json")
	return cleaned
}
