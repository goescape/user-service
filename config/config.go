package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	BaseURL  string
	Port     string
	Grpc     RPCConfig
	Postgres PostgreSQLConfig
	Redis    RedisConfig
	Kafka    KafkaConfig

	ServiceOrderAdress string
}

func Load() (*Config, error) {
	viper.New()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed read config file: %v", err)
	}

	cfg := &Config{
		BaseURL: viper.GetString("BASE_URL_PATH"),
		Port:    viper.GetString("PORT"),

		Grpc: RPCConfig{
			Port: viper.GetString("RPC_PORT"),
		},

		Postgres: PostgreSQLConfig{
			DbHost:        viper.GetString("DB_HOST"),
			DbPort:        viper.GetString("DB_PORT"),
			DbUsername:    viper.GetString("DB_USERNAME"),
			DbPassword:    viper.GetString("DB_PASSWORD"),
			DbName:        viper.GetString("DB_NAME"),
			DbMaxOpenConn: viper.GetInt("DB_MAX_OPEN_CONN"),
			DbMaxIdleConn: viper.GetInt("DB_MAX_IDLE_CONN"),
			DbMaxLifeTime: viper.GetDuration("DB_MAX_LIFE_TIME"),
		},

		Redis: RedisConfig{
			Address:  viper.GetString("REDIS_ADDRESS"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},

		Kafka: KafkaConfig{
			Address: []string{viper.GetString("KAFKA_ADDRESS")},
		},

		ServiceOrderAdress: viper.GetString("SERVICE_ORDER_ADDRESS"),
	}

	return cfg, nil
}
