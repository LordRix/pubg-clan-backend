package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pubg-clan-backend/models"
	"pubg-clan-backend/utils"
	"strconv"
	"strings"
	"time"
)

var (
	pubgBaseURL = "https://api.pubg.com/shards/steam"
	apiKey      string // initialized by InitAPIKey()
	httpClient  = &http.Client{Timeout: 10 * time.Second}
	ClanMembers = []string{"LordRix", "TanaX18", "M1key-D", "GRenes", "Jarm00725", "Remnorz", "Kalliyo"} // customize your players
)

// InitAPIKey loads and validates the API key from environment
func InitAPIKey() {
	apiKey = strings.TrimSpace(os.Getenv("PUBG_API_KEY"))
	if apiKey == "" {
		log.Fatal("[FATAL] PUBG_API_KEY is empty or missing. Cannot start server.")
	}
}

// GetScoreboard retrieves the chicken dinners per player
func GetScoreboard(minDate time.Time) []models.ScoreboardEntry {
	cache.Lock.RLock()
	noCache := true
	if time.Now().Before(cache.ExpiryTime) && len(cache.Data) > 0 && !noCache {
		log.Println(utils.Green("[CACHE] Serving from cache"))
		result := cache.Data
		cache.Lock.RUnlock()
		return result
	}
	cache.Lock.RUnlock()

	log.Println(utils.Yellow("[CACHE] Cache expired or empty, rebuilding..."))

	var scoreboard []models.ScoreboardEntry

	playerIDMap, err := LoadPlayerIDMap("players_id.json")
	if err != nil {
		log.Printf("%s Failed to load player ID map: %v", utils.Red("[ERROR]"), err)
		return nil
	}

	for playerName, playerID := range playerIDMap {
		log.Printf("Player Name: %s, Player ID: %s", playerName, playerID)
		matchIDs, err := getPlayerMatches(playerName)
		if err != nil {
			log.Printf("%s Failed fetching matches for %s: %v", utils.Red("[ERROR]"), playerName, err)
			continue
		}
		log.Printf("%s : %s : %s\n\n", playerName, playerID, matchIDs)
	}

	// for _, playerName := range ClanMembers {
	// 	log.Printf("%s %s", utils.Blue("[PLAYER] Processing"), playerName)

	// 	playerID, err := GetOrFetchPlayerID(playerName)
	// 	if err != nil {
	// 		log.Printf("%s Failed fetching player ID for %s: %v", utils.Red("[ERROR]"), playerName, err)
	// 		continue
	// 	}
	// 	log.Printf("%s Fetching PlayerID for %s - %s", utils.Blue("[PUBG API]"), utils.Yellow(playerName), utils.Green(playerID))

	// 	matchIDs, err := getPlayerMatches(playerName)
	// 	if err != nil {
	// 		log.Printf("%s Failed fetching matches for %s: %v", utils.Red("[ERROR]"), playerName, err)
	// 		continue
	// 	}

	// 	chickenDinners := 0
	// 	for _, matchID := range matchIDs {
	// 		won, matchTime, err := checkIfChickenDinner(playerID, matchID)
	// 		if err != nil {
	// 			log.Printf("%s Failed checking match %s for %s: %v", utils.Red("[ERROR]"), matchID, playerName, err)
	// 			continue
	// 		}
	// 		if matchTime.Before(minDate) {
	// 			continue
	// 		}
	// 		if won {
	// 			chickenDinners++
	// 			log.Printf("%s %s won match %s 🐔", utils.Green("[MATCH]"), playerName, matchID)
	// 		}
	// 	}

	// 	log.Printf("%s %s has %d Chicken Dinners 🐔", utils.Green("[SUMMARY]"), playerName, chickenDinners)

	// 	scoreboard = append(scoreboard, models.ScoreboardEntry{
	// 		PlayerName:     playerName,
	// 		ChickenDinners: chickenDinners,
	// 	})
	// }

	// cache.Lock.Lock()
	// cache.Data = scoreboard
	// cache.ExpiryTime = time.Now().Add(cache.TTL)
	// cache.Lock.Unlock()

	// log.Println(utils.Green("[CACHE] Scoreboard cached successfully"))
	return scoreboard
}

// GetOrFetchPlayerID retrieves or caches a player's ID
func GetOrFetchPlayerID(playerName string) (string, error) {
	return getPlayerID(playerName)
}

// getPlayerID queries PUBG API to find a player ID
func getPlayerID(playerName string) (string, error) {
	cache.PlayerIDCacheMux.RLock()
	id, found := cache.PlayerIDCache[playerName]
	cache.PlayerIDCacheMux.RUnlock()

	if found {
		log.Printf("%s Using cached PlayerID for %s", utils.Green("[CACHE]"), playerName)
		return id, nil
	}

	req, _ := http.NewRequest("GET", pubgBaseURL+"/players?filter[playerNames]="+playerName, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := safeDoRequest(req)
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

	id = pr.Data[0].ID

	cache.PlayerIDCacheMux.Lock()
	cache.PlayerIDCache[playerName] = id
	cache.PlayerIDCacheMux.Unlock()
	//log.Printf("%s Fetching PlayerID for %s - %s", utils.Blue("[PUBG API]"), utils.Yellow(playerName), utils.Green(id))
	return id, nil
}

// getPlayerMatches fetches recent matches for a player
func getPlayerMatches(playerName string) ([]string, error) {
	req, _ := http.NewRequest("GET", pubgBaseURL+"/players?filter[playerNames]="+playerName, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := safeDoRequest(req)
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

	var matches []string
	for _, m := range pr.Data[0].Relationships.Matches.Data {
		matches = append(matches, m.ID)
	}

	return matches, nil
}

// checkIfChickenDinner checks if player won a specific match
func checkIfChickenDinner(playerID, matchID string) (bool, time.Time, error) {
	req, _ := http.NewRequest("GET", pubgBaseURL+"/matches/"+matchID, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := safeDoRequest(req)
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
	for _, participant := range match.Included {
		if participant.Type == "participant" && participant.Attributes.Stats.PlayerId == playerID {
			winPlace := participant.Attributes.Stats.WinPlace
			log.Println(utils.Green("Place: " + strconv.Itoa(winPlace)))
			return winPlace == 1, match.Data.Attributes.CreatedAt, nil
		}
	}
	return false, match.Data.Attributes.CreatedAt, fmt.Errorf("participant not found")
}

// safeDoRequest retries if 429 is received
func safeDoRequest(req *http.Request) (*http.Response, error) {
	maxRetries := 3
	backoff := time.Second

	for i := 0; i < maxRetries; i++ {
		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 429 {
			return resp, nil
		}

		log.Println(utils.Yellow("[RETRY] 429 Rate limit hit. Backing off..."))
		resp.Body.Close()
		time.Sleep(backoff)
		backoff *= 2
	}

	return nil, fmt.Errorf("too many retries after rate limiting")
}
func LoadPlayerIDMap(filePath string) (map[string]string, error) {
	// Read the file contents directly
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Parse the JSON into a map
	playerIDMap := make(map[string]string)
	if err := json.Unmarshal(data, &playerIDMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return playerIDMap, nil
}
