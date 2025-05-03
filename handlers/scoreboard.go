package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"pubg-clan-backend/services"
	"pubg-clan-backend/utils"
	"time"
)

var MinDate time.Time

func ScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.Blue("[HTTP] /scoreboard request received"))
	w.Header().Set("Content-Type", "application/json")

	scoreboard := services.GetScoreboard(MinDate)

	json.NewEncoder(w).Encode(scoreboard)
	log.Println(utils.Green("[HTTP] /scoreboard response sent"))
}
