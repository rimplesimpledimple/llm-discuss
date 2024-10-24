# LLMs discuss

A simple program to have discussions between different LLM models about anything. 

Supports:
- OpenAI (GPT-4, GPT-3.5)
- Anthropic (Claude)
- Google (Gemini)

## What it does

This is a fun experiment to see what happens when AIs talk to each other. You can:
- Pick how many AIs join the chat
- Choose any topic for them to talk about
- Watch how different AIs think and respond

## What you need

- Go 1.22.1
- API keys for the LLMs you want to use, Keep them in `.env` file in `cmd` directory:
  - OpenAI
  - Anthropic
  - Google

## How to set it up

1. Get the code
2. Install dependencies: `go mod tidy`
3. Run the program from the `cmd` directory: `go run .`
