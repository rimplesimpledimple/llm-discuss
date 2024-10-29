# LLM discuss

A simple program to have discussions between different LLM models about anything.

Supports:

- OpenAI
- Anthropic (Claude)
- Google (Gemini)
- Ollama
- DeepSeek

## What it does

This is a fun experiment to see what happens when AIs talk to each other. You can:

- Pick how many AIs join the chat
- Choose any topic for them to talk about
- Watch how different AIs think and respond

## What you need

- Go 1.22.1
- API keys for the LLMs you want to use, Keep them in `.env` file in `cmd` directory:
  - OpenAI if you want to use OpenAI
  - Anthropic if you want to use Claude
  - Google if you want to use Gemini
  - Ollama host and model specified if you want to use Ollama

## How to set it up

1. Get the code: `git clone https://github.com/rimplesimpledimple/llm-discuss.git`
2. Install dependencies: `go mod tidy`
3. Update the `.env` to add keys / URLs for the LLMs you want to use, or edit the `cmd/main.go` to customise the providers.
4. Run the program from the `cmd` directory: `go run .`
