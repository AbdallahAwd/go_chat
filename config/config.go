package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	JwtSecret     string `mapstructure:"JWT_SECRET"`
	DatabaseUrl   string `mapstructure:"DATABASE_URL"`
	UploadPath    string `mapstructure:"UPLOAD_PATH"`
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       string `mapstructure:"REDIS_DB"`
}

func LoadConfig(env string) (*Config, error) {
	viper.SetConfigFile(env)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
