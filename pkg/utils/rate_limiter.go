package utils

import (
	"sync"

	"golang.org/x/time/rate"
)

var userLimiters = make(map[string]*rate.Limiter, 0)

var mu sync.Mutex

func GetUserLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	limiter, existed := userLimiters[ip]

	if !existed {
		limiter = rate.NewLimiter(1, 1)
		userLimiters[ip] = limiter
	}
	return limiter
}
