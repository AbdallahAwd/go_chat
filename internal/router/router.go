package router

import (
	"chat_app/config"
	middlewares "chat_app/pkg/middleware"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, client *redis.Client, cnf *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// **** File Server
	fs := http.FileServer(http.Dir(cnf.UploadPath))
	r.Handle("/"+cnf.UploadPath+"/*", http.StripPrefix("/"+cnf.UploadPath+"/", fs))
	// ***
	authHandler := InitializeAuth(db, cnf)
	r.Route("/v1", func(r chi.Router) {
		r.Use(middlewares.RateLimiter(client))
		r.Post("/validate", authHandler.ValidatePhone)
		r.Route("/", func(r chi.Router) {
			r.Use(middlewares.PhoneMiddleware(cnf.JwtSecret))
			r.Post("/verify", authHandler.VerifyPhone)
			r.Post("/login", authHandler.CreateUser)
		})
		r.Route("/user", func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(cnf.JwtSecret))
			r.Get("/", authHandler.GetUser)
			r.Get("/all", authHandler.GetAllUser)
		})
	})
	return r
}
