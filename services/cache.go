package services

import (
	"pubg-clan-backend/models"
	"pubg-clan-backend/utils"
	"sync"
	"time"
)

var cache = Cache{
	TTL:           utils.GetEnvDuration("CACHE_DURATION", 10*time.Minute),
	PlayerIDCache: make(map[string]string), // 🛠 Initialize map early
}

type Cache struct {
	Data             []models.ScoreboardEntry
	ExpiryTime       time.Time
	PlayerIDCache    map[string]string
	PlayerIDCacheMux sync.RWMutex
	Lock             sync.RWMutex
	TTL              time.Duration
}
