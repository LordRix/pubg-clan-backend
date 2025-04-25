package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"pubg-clan-backend/models"
	"time"
)

var (
	pubgBaseURL = "https://api.pubg.com/shards/steam"
	apiKey      = os.Getenv("PUBG_API_KEY")
	httpClient  = &http.Client{Timeout: 10 * time.Second}
	clanMembers = []string{"Friend1", "Friend2", "Friend3"}
)

func GetScoreboard(minDate time.Time) []models.ScoreboardEntry {
	// Use cache first
	cache.Lock.RLock()
	if time.Now().Before(cache.ExpiryTime) && len(cache.Data) > 0 {
		result := cache.Data
		cache.Lock.RUnlock()
		return result
	}
	cache.Lock.RUnlock()

	var scoreboard []models.ScoreboardEntry

	for _, playerName := range clanMembers {
		playerID, err := getPlayerID(playerName)
		if err != nil {
			continue
		}

		matchIDs, err := getPlayerMatches(playerName)
		if err != nil {
			continue
		}

		chickenDinners := 0
		for _, matchID := range matchIDs {
			won, matchTime, err := checkIfChickenDinner(playerID, matchID)
			if err != nil {
				continue
			}
			if matchTime.Before(minDate) {
				continue
			}
			if won {
				chickenDinners++
			}
		}

		scoreboard = append(scoreboard, models.ScoreboardEntry{
			PlayerName:     playerName,
			ChickenDinners: chickenDinners,
		})
	}

	// Update cache
	cache.Lock.Lock()
	cache.Data = scoreboard
	cache.ExpiryTime = time.Now().Add(cache.TTL)
	cache.Lock.Unlock()

	return scoreboard
}

func getPlayerID(playerName string) (string, error) {
	req, _ := http.NewRequest("GET", pubgBaseURL+"/players?filter[playerNames]="+playerName, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("PUBG API error: %s", resp.Status)
	}

	var pr models.PlayerResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return "", err
	}

	if len(pr.Data) == 0 {
		return "", fmt.Errorf("no player data found")
	}

	return pr.Data[0].ID, nil
}

func getPlayerMatches(playerName string) ([]string, error) {
	req, _ := http.NewRequest("GET", pubgBaseURL+"/players?filter[playerNames]="+playerName, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PUBG API error: %s", resp.Status)
	}

	var pr models.PlayerResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}

	if len(pr.Data) == 0 {
		return nil, fmt.Errorf("no player matches found")
	}

	var matches []string
	for _, m := range pr.Data[0].Relationships.Matches.Data {
		matches = append(matches, m.ID)
	}

	return matches, nil
}

func checkIfChickenDinner(playerID, matchID string) (bool, time.Time, error) {
	req, _ := http.NewRequest("GET", pubgBaseURL+"/matches/"+matchID, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, time.Now(), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, time.Now(), fmt.Errorf("PUBG API error: %s", resp.Status)
	}

	var match models.MatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&match); err != nil {
		return false, time.Now(), err
	}

	for _, participant := range match.Data.Included {
		if participant.Type == "participant" && participant.ID == playerID {
			winPlace := participant.Attributes.Stats.WinPlace
			return winPlace == 1, match.Data.Attributes.CreatedAt, nil
		}
	}

	return false, match.Data.Attributes.CreatedAt, fmt.Errorf("participant not found in match")
}
