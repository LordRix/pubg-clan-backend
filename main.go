package main

import (
	"log"
	"net/http"
	"pubg-clan-backend/handlers"
	"pubg-clan-backend/services"
	"pubg-clan-backend/utils"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(utils.Yellow("No .env file found, using system environment variables"))
	}

	services.InitAPIKey() // <-- Load and validate API Key early

	port := utils.GetEnv("PORT", "8080")
	minDateStr := utils.MustGetEnv("MIN_DATE")
	var err error
	handlers.MinDate, err = time.Parse(time.RFC3339, minDateStr)
	if err != nil {
		log.Fatalf(utils.Red("[INIT ERROR] Invalid MIN_DATE format: %v"), err)
	}

	log.Println(utils.Blue("[INIT] Warming up player cache..."))
	if err := handlers.WarmupPlayerCache(); err != nil {
		log.Fatalf(utils.Red("[INIT ERROR] Failed warmup: %v"), err)
	}

	http.HandleFunc("/scoreboard", handlers.ScoreboardHandler)

	log.Printf(utils.Blue("Server running on port %s...\n"), port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
