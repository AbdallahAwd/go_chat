package internal

import (
	"chat_app/config"
	"chat_app/internal/handlers"
	"chat_app/internal/repositories"
	"chat_app/internal/services"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeAuthHanlder(db *gorm.DB, conf *config.Config) *handlers.AuthHandler {
	wire.Build(
		repositories.NewAuthRepository,
		services.NewAuthService,
		handlers.NewAuthHandler,
	)
	return &handlers.AuthHandler{}
}
func InitializeChatHanlder(db *gorm.DB, conf *config.Config) *handlers.ChatHandler {
	wire.Build(
		repositories.NewChatRepo,
		services.NewChatService,
		handlers.NewChathandler,
	)
	return &handlers.ChatHandler{}
}
