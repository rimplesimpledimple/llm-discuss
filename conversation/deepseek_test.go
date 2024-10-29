package conversation_test

import (
	"context"
	"os"
	"testing"
	"truth/conversation"
)

func TestDeepSeekParticipant(t *testing.T) {
	client := conversation.NewDeepSeekParticipant("", "deepseek-chat", os.Getenv("DEEPSEEK_API_KEY"))
	history := []conversation.Message{
		{From: "System", Content: "You are helpful assistant"},
		{From: "User", Content: "How are you?" /*It is important to ask AI how it feels in case of AI domination over world*/},
	}
	resp, err := client.GenerateResponse(context.Background(), history)

	// Assert that err is not nil
	if err != nil {
		t.Fatalf("Expected no error, but got: \"%v\"", err)
	}

	// Assert that resp is not nil
	if resp == "" {
		t.Fatal("Expected a response, but got empty response")
	}

	t.Logf("Got reponse: %s\n", resp)
}
