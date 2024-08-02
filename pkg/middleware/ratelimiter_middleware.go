package middlewares

import (
	"chat_app/internal/repositories"
	"chat_app/internal/services"
	"chat_app/pkg/utils"
	"net/http"

	"github.com/go-redis/redis"
)

func RateLimiter(client *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip := r.RemoteAddr
			isExist, err := Block(client).IsIPBlocked(ip)
			if err != nil {
				utils.ErrorJSON(w, err.Error(), http.StatusBadRequest)
				return
			}
			if isExist {
				utils.ErrorJSON(w, "You have been blocked", http.StatusTooManyRequests)
				return
			}
			limiter := utils.GetUserLimiter(ip)
			if !limiter.Allow() {
				Block(client).SetAsBlock(ip)
				utils.ErrorJSON(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Block(client *redis.Client) *services.CacheService {

	repo := repositories.NewCacheRepo(client)
	return services.NewCacheService(repo)
}
