package config

import (
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	AppName     string `mapstructure:"APP_NAME"`
	Debug       bool   `mapstructure:"DEBUG"`

	PostgresDSN string `mapstructure:"DB_POSTGRES_DSN"`
	ResendKey   string `mapstructure:"RESEND_API_KEY"`

	JWTSecretKey            string        `mapstructure:"JWT_SECRET"`
	JWTAccessTokenDuration  time.Duration `mapstructure:"JWT_ACCESS_TOKEN_DURATION"`
	JWTRefreshTokenDuration time.Duration `mapstructure:"JWT_REFRESH_TOKEN_DURATION"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       string `mapstructure:"REDIS_DB"`
}

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("error reading config file, %s", err)
		}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("error unmarshalling config, %s", err)
		}
	})

	return config
}
