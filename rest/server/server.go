package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"truth/conversation"
)

type Configuration struct {
	InitialPrompt            string `json:"initial_prompt"`
	OpenAiParticipantSize    int    `json:"openai_participant_size"`
	AnthropicParticipantSize int    `json:"anthropic_participant_size"`
	DeepSeekParticipantSize  int    `json:"deepseek_participant_size"`
	GeminiParticipantSize    int    `json:"gemini_participant_size"`
	Turns                    int    `json:"turns"`
}

func NewConfiguration(initialPrompt string, openAiParticipantSize, anthropicParticipantSize, deepSeekParticipantSize, geminiParticipantSize, turns int) *Configuration {
	return &Configuration{
		InitialPrompt:            initialPrompt,
		OpenAiParticipantSize:    openAiParticipantSize,
		AnthropicParticipantSize: anthropicParticipantSize,
		DeepSeekParticipantSize:  deepSeekParticipantSize,
		GeminiParticipantSize:    geminiParticipantSize,
		Turns:                    turns,
	}
}

type Server struct {
	Config             *Configuration
	historyListeners   map[uint]chan conversation.Message
	historyListenersId uint

	participants []conversation.Participant
	conversation *conversation.Conversation
}

func NewServer(c *Configuration) *Server {
	return &Server{c, make(map[uint]chan conversation.Message), 0, make([]conversation.Participant, 0), nil}
}

func (s *Server) Run() {
	go func() {
		ctx := context.Background()
		s.addParcipicants()
		s.conversation = conversation.NewConversation(s.participants)
		err := s.conversation.Start(fmt.Sprintf(s.Config.InitialPrompt, len(s.Config.InitialPrompt)))
		if err != nil {
			log.Fatalf("Error starting conversation: %v", err)
		}
		for i := 0; i < s.Config.Turns; i++ {
			if err := s.conversation.NextTurn(ctx); err != nil {
				log.Printf("Error in conversation turn: %v", err)
				break
			}
			s.notifyHistoryListeners(s.conversation.History[len(s.conversation.History)-1])
		}
	}()
}

func (s *Server) addParcipicants() {
	userCount := 0
	// Create participants slice with non-nil values only
	var participants []conversation.Participant

	// Add OpenAI participant if API key is set
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		for i := 0; i < s.Config.OpenAiParticipantSize; i++ {
			userCount++
			participants = append(participants,
				conversation.NewOpenAIParticipant("User "+fmt.Sprint(userCount), "gpt-4", key))
		}
	}

	// Add Claude participant if API key is set
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		for i := 0; i < s.Config.AnthropicParticipantSize; i++ {
			userCount++
			participants = append(participants,
				conversation.NewClaudeParticipant("User "+fmt.Sprint(userCount), "claude-3-5-sonnet-20240620", key))
		}
	}

	// Add Gemini participant if API key is set
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		for i := 0; i < s.Config.GeminiParticipantSize; i++ {
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
		for i := 0; i < s.Config.DeepSeekParticipantSize; i++ {
			userCount++
			participants = append(participants,
				conversation.NewDeepSeekParticipant("User "+fmt.Sprint(userCount), "deepseek-chat", key))
		}
	}

	s.participants = participants
}

func (s *Server) notifyHistoryListeners(msg conversation.Message) {
	for _, channel := range s.historyListeners {
		channel <- msg
	}
}

// returns (id, chan)
// first agrumnet indicates if we want to send all messages(one that happend in past too) or only one that happens after listener
func (s *Server) AddHistoryListener(all bool) (uint, <-chan conversation.Message) {
	ch := make(chan conversation.Message)
	id := s.historyListenersId
	s.historyListenersId++
	s.historyListeners[id] = ch
	if all {
		go func() {
			if s.conversation != nil {
				for i := range s.conversation.History {
					ch <- s.conversation.History[i]
				}
			}
		}()
	}
	return id, ch
}

func (s *Server) RemoveHistoryListener(id uint) {
	close(s.historyListeners[id])
	delete(s.historyListeners, id)
}
