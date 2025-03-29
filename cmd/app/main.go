package main

import (
	"expat-news/queue-manager/internal/handlers/msg"
	"expat-news/queue-manager/internal/handlers/publisher"
	"expat-news/queue-manager/internal/handlers/queue"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/msg", msg.Handler)
	http.HandleFunc("/publisher", publisher.Handler)
	http.HandleFunc("/queue", queue.Handler)
	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
