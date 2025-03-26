package main

import (
	"expat-news/queue-manager/internal/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/msg", msg.Handler)
	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
