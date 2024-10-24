package conversation

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiParticipant struct {
	name   string
	config *ModelConfig
	model_ *genai.GenerativeModel
}

func NewGeminiParticipant(
	name string,
	model string,
	apiKey string,
) *GeminiParticipant {
	ctx := context.Background()
	config, ok := GetModelConfig(model)
	if !ok {
		log.Fatalf("unknown model: %s", model)
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed to create Gemini client: %v", err)
	}

	return &GeminiParticipant{
		name:   name,
		config: &config,
		model_: client.GenerativeModel(model),
	}
}

func (p *GeminiParticipant) GetName() string {
	return fmt.Sprintf("%s (%s)", p.name, p.config.Name)
}

func (p *GeminiParticipant) GenerateResponse(ctx context.Context, history []Message) (string, error) {
	chat := p.model_.StartChat()

	// Convert our Message format to Gemini's format
	for _, msg := range history {
		if msg.From == "System" {
			systemPrompt := fmt.Sprintf("You are %s. %s. Respond in format: %s: <response>",
				p.GetName(), msg.Content, p.GetName())
			// For Gemini, we'll send system message as a user message since it doesn't have a dedicated system role
			p.model_.SystemInstruction = &genai.Content{
				Parts: []genai.Part{
					genai.Text(systemPrompt),
				},
				Role: "user",
			}
		} else if msg.From == p.GetName() {
			chat.History = append(chat.History, &genai.Content{
				Parts: []genai.Part{
					genai.Text(msg.Content),
				},
				Role: "model",
			})
		} else {
			chat.History = append(chat.History, &genai.Content{
				Parts: []genai.Part{
					genai.Text(msg.Content),
				},
				Role: "user",
			})
		}
	}

	// Manage conversation context window
	chat.History = p.manageConversation(chat.History)

	// Generate response
	resp, err := chat.SendMessage(ctx, genai.Text("Continue the conversation"))
	if err != nil {
		return "", fmt.Errorf("error generating response: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	// Extract the text from the response
	return string(resp.Candidates[0].Content.Parts[0].(genai.Text)), nil
}

func (p *GeminiParticipant) manageConversation(history []*genai.Content) []*genai.Content {
	currentTokens := p.countTokens(history)
	// Reserve tokens for the response
	for currentTokens+p.config.MaxTokens > p.config.ContextWindow {
		if len(history) <= 1 { // Keep at least one message
			break
		}
		// Remove the second message
		history = append(history[:1], history[2:]...)
		currentTokens = p.countTokens(history)
	}
	return history
}

func (p *GeminiParticipant) countTokens(history []*genai.Content) int {
	// Simple estimation: ~4 chars per token
	totalChars := 0
	for _, content := range history {
		for _, part := range content.Parts {
			totalChars += len(part.(genai.Text))
		}
	}
	return totalChars / 4
}
