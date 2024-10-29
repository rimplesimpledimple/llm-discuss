package conversation

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Message represents a single message in the conversation
type Message struct {
	From      string    `json:"from"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Participant interface for different AI models
type Participant interface {
	GetName() string
	GenerateResponse(ctx context.Context, history []Message) (string, error)
}

// Conversation manages the interaction between participants
type Conversation struct {
	Participants []Participant
	History      []Message
	CurrentTurn  int
}

// NewConversation creates a new conversation with the given participants
func NewConversation(participants []Participant) *Conversation {
	return &Conversation{
		Participants: participants,
		History:      make([]Message, 0),
		CurrentTurn:  0,
	}
}

// Start begins the conversation with an initial prompt
func (c *Conversation) Start(initialPrompt string) error {
	c.History = append(c.History, Message{
		From:      "System",
		Content:   initialPrompt,
		Timestamp: time.Now(),
	})
	return nil
}

// NextTurn advances the conversation to the next participant
func (c *Conversation) NextTurn(ctx context.Context) error {
	participant := c.Participants[c.CurrentTurn%len(c.Participants)]

	response, err := participant.GenerateResponse(ctx, c.History)
	if err != nil {
		return fmt.Errorf("error generating response: %v", err)
	}

	// Get color for the current participant
	color := GetParticipantColor(c.CurrentTurn % len(c.Participants))
	// Print the colored log
	log.Printf("%s%s%s\n", color, response, ColorReset)

	c.History = append(c.History, Message{
		From:      participant.GetName(),
		Content:   response,
		Timestamp: time.Now(),
	})

	c.CurrentTurn++
	return nil
}

// PrintHistory displays the entire conversation with colors
func (c *Conversation) PrintHistory() {
	participantColors := make(map[string]string)

	// Assign colors to participants
	for i, p := range c.Participants {
		participantColors[p.GetName()] = GetParticipantColor(i)
	}

	// Special color for system messages
	participantColors["System"] = ColorCyan

	for _, msg := range c.History {
		color := participantColors[msg.From]
		if color == "" {
			color = ColorReset
		}

		fmt.Printf("%s[%s] %s: %s%s\n",
			color,
			msg.Timestamp.Format("15:04:05"),
			msg.From,
			msg.Content,
			ColorReset,
		)
	}
}
