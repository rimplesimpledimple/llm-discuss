package main

import (
	"fmt"
	"net/http"
	"truth/rest/handlers"
	"truth/rest/server"
)

const (
	initialPrompt = "This is a discussion between %d participants." +
		"Your mission is to discuss whether AI should be closed source or open source. Consider the implications for innovation, " +
		"safety, transparency, and societal impact. Discuss the pros and cons of both approaches, potential hybrid models, " +
		"and the role of regulation in AI development. You are one of the participants. Write short anwsers."
)

func main() {
	config := loadConfig()
	configuration := server.NewConfiguration(initialPrompt, 0, 0, 5, 0, 100)
	configurationHandler := handlers.NewConfigurationHandler(configuration)

	server := server.NewServer(configuration)
	startHandler := handlers.NewStartHandler(server)

	historyHandler := handlers.NewHistoryHandler(server)

	http.Handle("/configure", configurationHandler)
	http.Handle("/start", startHandler)
	http.Handle("/history/stream", historyHandler)

	fmt.Printf("Starting server on :%d\n", config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
