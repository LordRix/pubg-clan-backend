package models

import "time"

type PlayerResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Name string `json:"name"`
		} `json:"attributes"`
		Relationships struct {
			Matches struct {
				Data []struct {
					ID string `json:"id"`
				} `json:"data"`
			} `json:"relationships"`
		} `json:"relationships"`
	} `json:"data"`
}

type MatchResponse struct {
	Data struct {
		ID         string `json:"id"`
		Attributes struct {
			CreatedAt time.Time `json:"createdAt"`
		} `json:"attributes"`
		Included []struct {
			Type       string `json:"type"`
			ID         string `json:"id"`
			Attributes struct {
				Stats struct {
					WinPlace int `json:"winPlace"`
				} `json:"stats"`
			} `json:"attributes"`
		} `json:"included"`
	} `json:"data"`
}

type ScoreboardEntry struct {
	PlayerName     string `json:"playerName"`
	ChickenDinners int    `json:"chickenDinners"`
}
