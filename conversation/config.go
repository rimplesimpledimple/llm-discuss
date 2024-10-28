package conversation

import (
	"encoding/json"
	"log"
	"os"
)

// ModelConfig holds configuration for each model
type ModelConfig struct {
	Name          string
	ContextWindow int
	MaxTokens     int
	Temperature   float32
}

// ModelConfigs maps model names to their configurations
var ModelConfigs = map[string]ModelConfig{
	"claude-3-5-sonnet-20241022-v2": {
		Name:          "claude-3-5-sonnet-20241022-v2",
		ContextWindow: 128000,
		MaxTokens:     4096,
		Temperature:   0.7,
	},
	"gemini-1.5-flash": {
		Name:          "gemini-1.5-flash",
		ContextWindow: 128000,
		MaxTokens:     2048,
		Temperature:   0.7,
	},
	"gpt-4": {
		Name:          "gpt-4",
		ContextWindow: 128000,
		MaxTokens:     2048,
		Temperature:   0.7,
	},
	"qwen2.5:7b-instruct-q6_K": {
		Name:          "qwen2.5:7b-instruct-q6_K",
		ContextWindow: 4096,
		MaxTokens:     2048,
		Temperature:   0.7,
	},
	"deepseek-chat": {
		Name:          "deepseek-chat",
		ContextWindow: 4096,
		MaxTokens:     2048,
		Temperature:   0.7,
	},
}

// GetModelConfig returns the configuration for a given model
func GetModelConfig(model string) (ModelConfig, bool) {
	config, exists := ModelConfigs[model]
	return config, exists
}

// Read OLLAMA_MODEL_CONFIG from env, if it exists validate it and update ModelConfigs
func OllamaModelConfig() {
	ollamaModelConfig := os.Getenv("OLLAMA_MODEL_CONFIG")
	if ollamaModelConfig != "" {
		var modelConfigs map[string]ModelConfig
		err := json.Unmarshal([]byte(ollamaModelConfig), &modelConfigs)
		if err != nil {
			log.Fatalf("failed to unmarshal OLLAMA_MODEL_CONFIG: %v", err)
		}
		for model, config := range modelConfigs {
			if _, exists := ModelConfigs[model]; !exists {
				log.Fatalf("unknown model in OLLAMA_MODEL_CONFIG: %s", model)
			}
			ModelConfigs[model] = config
		}
	}
}
