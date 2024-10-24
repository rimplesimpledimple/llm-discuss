package conversation

import (
	"context"
	"fmt"
	"log"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type ClaudeParticipant struct {
	name   string
	model  string
	config *ModelConfig
	client *anthropic.Client
}

func NewClaudeParticipant(
	name string,
	model string,
	apiKey string,
) *ClaudeParticipant {
	config, ok := GetModelConfig(model)
	if !ok {
		log.Fatalf("unknown model: %s", model)
	}
	return &ClaudeParticipant{
		name:   name,
		config: &config,
		client: anthropic.NewClient(
			option.WithAPIKey(apiKey),
		),
	}
}

func (p *ClaudeParticipant) GetName() string {
	return fmt.Sprintf("%s (%s)", p.name, p.model)
}

func (p *ClaudeParticipant) GenerateResponse(ctx context.Context, history []Message) (string, error) {
	messages := make([]anthropic.MessageParam, 0, len(history))
	systemPrompt := ""
	// Convert our Message format to Anthropic's MessageParam format
	for _, msg := range history {
		if msg.From == "System" {
			systemPrompt = fmt.Sprintf("You are %s. %s. Respond in format: %s: <response>",
				p.GetName(), msg.Content, p.GetName())
		} else if msg.From == p.GetName() {
			messages = append(messages, anthropic.NewAssistantMessage(
				anthropic.NewTextBlock(msg.Content),
			))
		} else {
			messages = append(messages, anthropic.NewUserMessage(
				anthropic.NewTextBlock(msg.Content),
			))
		}
	}

	// Manage conversation context window
	messages = p.manageConversation(messages)

	resp, err := p.client.Messages.New(
		ctx,
		anthropic.MessageNewParams{
			Model:       anthropic.F(p.config.Name),
			Messages:    anthropic.F(messages),
			MaxTokens:   anthropic.F(int64(p.config.MaxTokens)),
			System:      anthropic.F([]anthropic.TextBlockParam{anthropic.NewTextBlock(systemPrompt)}),
			Temperature: anthropic.F(float64(p.config.Temperature)),
		},
	)

	if err != nil {
		return "", fmt.Errorf("error creating message: %v", err)
	}

	return resp.Content[0].Text, nil
}

func (p *ClaudeParticipant) manageConversation(messages []anthropic.MessageParam) []anthropic.MessageParam {
	currentTokens := p.countTokens(messages)
	// Reserve tokens for the response
	for currentTokens+p.config.MaxTokens > p.config.ContextWindow {
		if len(messages) <= 1 { // Keep at least one message
			break
		}
		// Remove the second message
		messages = append(messages[:1], messages[2:]...)
		currentTokens = p.countTokens(messages)
	}
	return messages
}

func (p *ClaudeParticipant) countTokens(messages []anthropic.MessageParam) int {
	// Simple estimation: ~4 chars per token
	totalChars := 0
	for _, msg := range messages {
		for _, content := range msg.Content.Value {
			if text, ok := content.(anthropic.TextBlockParam); ok {
				totalChars += len(text.Text.Value)
			}
		}
	}
	return totalChars / 4
}
