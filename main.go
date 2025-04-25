package main

import (
	"log"
	"net/http"
	"time"

	"pubg-clan-backend/handlers"
	"pubg-clan-backend/utils"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := utils.GetEnv("PORT", "8080")
	minDateStr := utils.MustGetEnv("MIN_DATE")
	var err error
	handlers.MinDate, err = time.Parse(time.RFC3339, minDateStr)
	if err != nil {
		log.Fatalf("Invalid MIN_DATE format: %v", err)
	}

	http.HandleFunc("/scoreboard", handlers.ScoreboardHandler)

	log.Printf("Server running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
