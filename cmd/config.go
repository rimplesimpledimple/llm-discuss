package main

type Config struct {
	OpenAiParticipantSize    int `json:"openai_participant_size"`
	AnthropicParticipantSize int `json:"anthropic_participant_size"`
	DeepSeekParticipantSize  int `json:"deepseek_participant_size"`
	GeminiParticipantSize    int `json:"gemini_participant_size"`
}

func NewDefaultConfig() *Config {
	return &Config{
		1, 1, 1, 1,
	}
}
