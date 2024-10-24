package main

import (
	"context"
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

const (
	initialPrompt = "This is a discussion between %d participants." +
		"Your mission is to discuss whether AI should be closed source or open source. Consider the implications for innovation, " +
		"safety, transparency, and societal impact. Discuss the pros and cons of both approaches, potential hybrid models, " +
		"and the role of regulation in AI development. "
)

func main() {
	loadEnv()
	ctx := context.Background()

	participants := []conversation.Participant{
		conversation.NewOpenAIParticipant("User 1", "gpt-4", os.Getenv("OPENAI_API_KEY")),
		conversation.NewClaudeParticipant("User 2", "claude-3-5-sonnet-20240620", os.Getenv("ANTHROPIC_API_KEY")),
		conversation.NewGeminiParticipant("User 3", "gemini-1.5-flash", os.Getenv("GEMINI_API_KEY")),
	}

	// Create a new conversation with multiple participants
	conversation := conversation.NewConversation(participants)

	// Start the conversation with an initial prompt
	err := conversation.Start(fmt.Sprintf(initialPrompt, len(conversation.Participants)))
	if err != nil {
		log.Fatalf("Error starting conversation: %v", err)
	}

	// Run for 5 turns
	for i := 0; i < 5; i++ {
		if err := conversation.NextTurn(ctx); err != nil {
			log.Printf("Error in conversation turn: %v", err)
			break
		}
	}
}
