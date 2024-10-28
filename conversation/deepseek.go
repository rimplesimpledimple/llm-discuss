package conversation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type DeepSeekParticipant struct {
	name   string
	apiKey string
	config *ModelConfig
}

func NewDeepSeekParticipant(
	name string,
	model string,
	apiKey string,
) *DeepSeekParticipant {
	config, ok := GetModelConfig(model)
	if !ok {
		log.Fatalf("unknown model: %s", model)
	}
	return &DeepSeekParticipant{
		name:   name,
		config: &config,
		apiKey: apiKey,
	}
}

func (p *DeepSeekParticipant) GetName() string {
	return fmt.Sprintf("%s (%s)", p.name, p.config.Name)
}

type DeepSeekChatCompletion struct {
	ID                string           `json:"id"`
	Choices           []DeepSeekChoice `json:"choices"`
	Created           int              `json:"created"`
	Model             string           `json:"model"`
	SystemFingerprint string           `json:"system_fingerprint"`
	Object            string           `json:"object"`
	Usage             DeepSeekUsage    `json:"usage"`
}

type DeepSeekChoice struct {
	FinishReason string           `json:"finish_reason"`
	Index        int              `json:"index"`
	Message      DeepSeekMessage  `json:"message"`
	Logprobs     DeepSeekLogprobs `json:"logprobs"`
}

type DeepSeekMessage struct {
	Content   string             `json:"content"`
	ToolCalls []DeepSeekToolCall `json:"tool_calls"`
	Role      string             `json:"role"`
}

type DeepSeekToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function DeepSeekFunction `json:"function"`
}

type DeepSeekFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type DeepSeekLogprobs struct {
	Content []DeepSeekLogprobContent `json:"content"`
}

type DeepSeekLogprobContent struct {
	Token       string                   `json:"token"`
	Logprob     float64                  `json:"logprob"`
	Bytes       []byte                   `json:"bytes"`
	TopLogprobs []DeepSeekTopLogprobItem `json:"top_logprobs"`
}

type DeepSeekTopLogprobItem struct {
	Token   string  `json:"token"`
	Logprob float64 `json:"logprob"`
	Bytes   []byte  `json:"bytes"`
}

type DeepSeekUsage struct {
	CompletionTokens      int `json:"completion_tokens"`
	PromptTokens          int `json:"prompt_tokens"`
	PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
	TotalTokens           int `json:"total_tokens"`
}

func (p *DeepSeekParticipant) GenerateResponse(ctx context.Context, history []Message) (string, error) {
	type messagePayload struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	var messagesPayload []messagePayload
	for i := range history {
		var role string
		if strings.Compare(strings.ToLower(history[i].From), "system") == 0 {
			role = "system"
		} else if strings.Compare(strings.ToLower(history[i].From), strings.ToLower(p.GetName())) == 0 {
			role = "assistant"
		} else {
			role = "user"
		}
		messagesPayload = append(messagesPayload, messagePayload{Role: role, Content: history[i].Content})
	}
	payloadMap := map[string]interface{}{
		"messages":          messagesPayload,
		"model":             "deepseek-chat",
		"frequency_penalty": 0,
		"max_tokens":        2048,
		"presence_penalty":  0,
		"response_format": map[string]interface{}{
			"type": "text",
		},
		"stop":           nil,
		"stream":         false,
		"stream_options": nil,
		"temperature":    1,
		"top_p":          1,
		"tools":          nil,
		"tool_choice":    "none",
		"logprobs":       false,
		"top_logprobs":   nil,
	}
	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var chatCompletion DeepSeekChatCompletion
	err = json.Unmarshal([]byte(body), &chatCompletion)
	if err != nil {
		return "", fmt.Errorf("error while unrmasaling json: %w, body: %s", err, string(body))
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (p *DeepSeekParticipant) manageConversation(messages []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
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

func (p *DeepSeekParticipant) countTokens(messages []openai.ChatCompletionMessage) int {
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
