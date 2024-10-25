package conversation

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/teilomillet/gollm"
)

type OllamaParticipant struct {
	name   string
	llm    gollm.LLM
	config *ModelConfig
}

func NewOllamaParticipant(
	name string,
	model string,
	endpoint string,
) *OllamaParticipant {
	config, ok := GetModelConfig(model)
	if !ok {
		log.Fatalf("unknown model: %s", model)
	}

	llm, err := gollm.NewLLM(
		gollm.SetProvider("ollama"),
		gollm.SetModel(model),
		gollm.SetLogLevel(gollm.LogLevelInfo),
		gollm.SetOllamaEndpoint(endpoint),
	)
	if err != nil {
		log.Fatalf("failed to create Ollama LLM: %v", err)
	}

	return &OllamaParticipant{
		name:   name,
		config: &config,
		llm:    llm,
	}
}

func (p *OllamaParticipant) GetName() string {
	return fmt.Sprintf("%s (%s)", p.name, p.config.Name)
}

func (p *OllamaParticipant) GenerateResponse(ctx context.Context, history []Message) (string, error) {
	// Build conversation history as a prompt
	var promptBuilder strings.Builder

	// Add system message first if it exists
	for _, msg := range history {
		if msg.From == "System" {
			promptBuilder.WriteString(fmt.Sprintf("System: You are %s. %s. Respond in format: %s: <response>\n\n",
				p.GetName(), msg.Content, p.GetName()))
			break
		}
	}

	// Add other messages in chronological order
	for _, msg := range history {
		if msg.From != "System" {
			promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", msg.From, msg.Content))
		}
	}

	// Create the prompt
	prompt := gollm.NewPrompt(promptBuilder.String())

	// Generate response
	response, err := p.llm.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating response: %v", err)
	}

	// Clean up response if needed
	response = strings.TrimSpace(response)

	// Ensure response format matches expected pattern
	if !strings.HasPrefix(response, p.GetName()) {
		response = fmt.Sprintf("%s: %s", p.GetName(), response)
	}

	return response, nil
}
