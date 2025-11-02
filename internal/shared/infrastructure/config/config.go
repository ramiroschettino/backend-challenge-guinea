package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env      string         
	Port     string         
	Database DatabaseConfig 
	RabbitMQ RabbitMQConfig 
	Log      LogConfig      
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type RabbitMQConfig struct {
	URL          string 
	Exchange     string 
	QueueUsers   string 
	PrefetchCount int   
}

type LogConfig struct {
	Level  string 
	Format string 
}

// Aca uso viper como pedia el pdf

func Load() (*Config, error) {

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	
	viper.SetDefault("ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("RABBITMQ_EXCHANGE", "backend_events")
	viper.SetDefault("RABBITMQ_QUEUE_USERS", "users_commands")
	viper.SetDefault("RABBITMQ_PREFETCH_COUNT", 10)

	_ = viper.ReadInConfig()

	return &Config{
		Env:  viper.GetString("ENV"),
		Port: viper.GetString("PORT"),
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			Name:     viper.GetString("DB_NAME"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:          viper.GetString("RABBITMQ_URL"),
			Exchange:     viper.GetString("RABBITMQ_EXCHANGE"),
			QueueUsers:   viper.GetString("RABBITMQ_QUEUE_USERS"),
			PrefetchCount: viper.GetInt("RABBITMQ_PREFETCH_COUNT"),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
	}, nil
}