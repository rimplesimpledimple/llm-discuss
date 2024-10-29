package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"truth/conversation"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func loadConfig() Config {
	// Define the file path for the configuration file
	configFilePath := "config.json"

	// Check if the configuration file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// If the file does not exist, create it using the default configuration
		defaultConfig := NewDefaultConfig()
		configJSON, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			log.Fatalf("Error marshalling default config: %v", err)
		}

		// Write the default configuration to the file
		err = os.WriteFile(configFilePath, configJSON, 0644)
		if err != nil {
			log.Fatalf("Error writing default config to file: %v", err)
		}

		return *defaultConfig
	}

	// If the file exists, load the configuration from it
	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	return config
}

const (
	initialPrompt = "This is a discussion between %d participants." +
		"Your mission is to discuss whether AI should be closed source or open source. Consider the implications for innovation, " +
		"safety, transparency, and societal impact. Discuss the pros and cons of both approaches, potential hybrid models, " +
		"and the role of regulation in AI development. "
)

func main() {
	loadEnv()
	ctx := context.Background()

	userCount := 0

	config := loadConfig()

	// Create participants slice with non-nil values only
	var participants []conversation.Participant

	// Add OpenAI participant if API key is set
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		for i := 0; i < config.OpenAiParticipantSize; i++ {
			userCount++
			participants = append(participants,
				conversation.NewOpenAIParticipant("User "+fmt.Sprint(userCount), "gpt-4", key))
		}
	}

	// Add Claude participant if API key is set
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		for i := 0; i < config.AnthropicParticipantSize; i++ {
			userCount++
			participants = append(participants,
				conversation.NewClaudeParticipant("User "+fmt.Sprint(userCount), "claude-3-5-sonnet-20240620", key))
		}
	}

	// Add Gemini participant if API key is set
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		for i := 0; i < config.GeminiParticipantSize; i++ {
			userCount++
			participants = append(participants,
				conversation.NewGeminiParticipant("User "+fmt.Sprint(userCount), "gemini-1.5-flash", key))
		}
	}

	// Add Ollama participant if host and model are set
	if host := os.Getenv("OLLAMA_HOST"); host != "" {
		if model := os.Getenv("OLLAMA_MODEL"); model != "" {
			userCount++
			participants = append(participants,
				conversation.NewOllamaParticipant("User "+fmt.Sprint(userCount), model, host))
		}
	}

	if key := os.Getenv("DEEPSEEK_API_KEY"); key != "" {
		for i := 0; i < config.DeepSeekParticipantSize; i++ {
			userCount++
			participants = append(participants,
				conversation.NewDeepSeekParticipant("User "+fmt.Sprint(userCount), "deepseek-chat", key))
		}
	}

	// Check if we have any participants
	if len(participants) == 0 {
		log.Fatal("No participants available. Please set at least one provider's environment variables.")
	}

	// Create a new conversation with multiple participants
	conversation := conversation.NewConversation(participants)

	// Start the conversation with an initial prompt
	err := conversation.Start(fmt.Sprintf(initialPrompt, len(participants)))
	if err != nil {
		log.Fatalf("Error starting conversation: %v", err)
	}

	turns := 5
	if t := os.Getenv("TURNS"); t != "" {
		turns, _ = fmt.Sscanf(t, "%d", &turns)
	}
	for i := 0; i < turns; i++ {
		if err := conversation.NextTurn(ctx); err != nil {
			log.Printf("Error in conversation turn: %v", err)
			break
		}
	}
}
