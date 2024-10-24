package conversation

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type OpenAIParticipant struct {
	name   string
	client *openai.Client
	config *ModelConfig
}

func NewOpenAIParticipant(
	name string,
	model string,
	apiKey string,
) *OpenAIParticipant {
	config, ok := GetModelConfig(model)
	if !ok {
		log.Fatalf("unknown model: %s", model)
	}
	return &OpenAIParticipant{
		name:   name,
		config: &config,
		client: openai.NewClient(apiKey),
	}
}

func (p *OpenAIParticipant) GetName() string {
	return fmt.Sprintf("%s (%s)", p.name, p.config.Name)
}

func (p *OpenAIParticipant) GenerateResponse(ctx context.Context, history []Message) (string, error) {
	messages := make([]openai.ChatCompletionMessage, 0, len(history))
	for _, msg := range history {
		role := openai.ChatMessageRoleUser
		if msg.From == "System" {
			role = openai.ChatMessageRoleSystem
			msg.Content = fmt.Sprintf("You are %s. %s. Respond in format: %s: <response>", p.GetName(), msg.Content, p.GetName())
		} else if msg.From == p.GetName() {
			role = openai.ChatMessageRoleAssistant
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Manage conversation context window
	messages = p.manageConversation(messages)

	resp, err := p.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       p.config.Name,
			Messages:    messages,
			Temperature: p.config.Temperature,
			MaxTokens:   p.config.MaxTokens,
			Stream:      false,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error creating chat completion: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

func (p *OpenAIParticipant) manageConversation(messages []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	currentTokens := p.countTokens(messages)
	// Reserve tokens for the response
	for currentTokens+p.config.MaxTokens > p.config.ContextWindow {
		if len(messages) <= 1 { // Keep at least the system message if present
			break
		}
		// Remove the second message (preserve system message if it exists)
		messages = append(messages[:1], messages[2:]...)
		currentTokens = p.countTokens(messages)
	}
	return messages
}

func (p *OpenAIParticipant) countTokens(messages []openai.ChatCompletionMessage) int {
	// GPT-4 and GPT-3.5 use cl100k_base tokenizer
	// Each message has a base overhead of ~4 tokens
	// Role names are also tokenized
	baseTokens := len(messages) * 4

	for _, msg := range messages {
		// Count role tokens (approximately)
		baseTokens += len(msg.Role) / 4

		// Count content tokens (approximately)
		// This is still a rough estimation, but better than simple character count
		words := len(strings.Fields(msg.Content))
		// Average English word is about 1.3 tokens
		baseTokens += int(float64(words) * 1.3)
	}

	return baseTokens
}
