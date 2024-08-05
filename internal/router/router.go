package router

import (
	"chat_app/config"

	"chat_app/pkg/analytics"
	middlewares "chat_app/pkg/middleware"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, client *redis.Client, cnf *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// **** File Server
	fs := http.FileServer(http.Dir(cnf.UploadPath))
	r.Handle("/"+cnf.UploadPath+"/*", http.StripPrefix("/"+cnf.UploadPath+"/", fs))
	// ***
	authHandler := InitializeAuthHanlder(db, cnf)
	chatHandler := InitializeChatHanlder(db, cnf)
	analyze := analytics.RunAnalyze()

	analyze.Init()
	c := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"PUT", "POST"},
			AllowedHeaders:   []string{"Origin", "Content-Type"},
			AllowCredentials: true},
	)
	// r.route
	r.HandleFunc("/ws", chatHandler.ChatWebSocket)

	r.Route("/v1", func(r chi.Router) {
		r.Use(c.Handler)
		// r.Use(middlewares.RateLimiter(client))
		r.Post("/validate", authHandler.ValidatePhone)
		r.Route("/", func(r chi.Router) {
			r.Use(middlewares.PhoneMiddleware(cnf.JwtSecret))
			r.Post("/verify", authHandler.VerifyPhone)
			r.Post("/login", authHandler.CreateUser)
		})
		r.Route("/user", func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(cnf.JwtSecret, analyze))
			r.Get("/", authHandler.GetUser)
			r.Get("/partners", chatHandler.GetMessagedUsers)
			r.Get("/messages", chatHandler.GetChatBetweenTwoUsers)
			r.Get("/all", authHandler.GetAllUser)
		})
		r.Route("/admin", func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(cnf.JwtSecret, analyze))
			r.Get("/metrics", promhttp.Handler().(http.HandlerFunc))

		})
	})
	return r
}
