package db

import (
	"chat_app/config"
	"chat_app/internal/models"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBRunner struct {
	DB     *gorm.DB
	Client *redis.Client
}

func RunDB(config *config.Config) (*DBRunner, error) {
	db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}
	postgressDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// conntection pooling
	postgressDB.SetMaxIdleConns(10)
	postgressDB.SetMaxOpenConns(100)
	postgressDB.SetConnMaxLifetime(time.Hour)
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       0,
	})

	return &DBRunner{DB: db, Client: client}, err
}

func (rDB *DBRunner) Migrate() error {
	return rDB.DB.AutoMigrate(&models.User{}, &models.Message{})
}
