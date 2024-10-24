package conversation

// ModelConfig holds configuration for each model
type ModelConfig struct {
	Name          string
	ContextWindow int
	MaxTokens     int
	Temperature   float32
}

// ModelConfigs maps model names to their configurations
var ModelConfigs = map[string]ModelConfig{
	"claude-3-5-sonnet-20240620": {
		Name:          "claude-3-5-sonnet-20240620",
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
}

// GetModelConfig returns the configuration for a given model
func GetModelConfig(model string) (ModelConfig, bool) {
	config, exists := ModelConfigs[model]
	return config, exists
}
