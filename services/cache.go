package services

import (
	"pubg-clan-backend/models"
	"pubg-clan-backend/utils"
	"sync"
	"time"
)

var cache = Cache{
	TTL: utils.GetEnvDuration("CACHE_DURATION", 10*time.Minute),
}

type Cache struct {
	Data       []models.ScoreboardEntry
	ExpiryTime time.Time
	Lock       sync.RWMutex
	TTL        time.Duration
}
