package handlers

import (
	"encoding/json"
	"net/http"
	"pubg-clan-backend/services"
	"time"
)

var MinDate time.Time

func ScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	scoreboard := services.GetScoreboard(MinDate)

	json.NewEncoder(w).Encode(scoreboard)
}
