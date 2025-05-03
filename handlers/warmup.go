package handlers

import (
	"pubg-clan-backend/services"
)

func WarmupPlayerCache() error {
	for _, playerName := range services.ClanMembers {
		_, err := services.GetOrFetchPlayerID(playerName)
		if err != nil {
			return err
		}
	}
	return nil
}
