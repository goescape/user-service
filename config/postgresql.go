package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type PostgreSQLConfig struct {
	DbHost        string
	DbPort        string
	DbUsername    string
	DbPassword    string
	DbName        string
	DbMaxOpenConn int
	DbMaxIdleConn int
	DbMaxLifeTime time.Duration
}

func InitPostgreSQL(cfg PostgreSQLConfig) (*sql.DB, error) {
	conn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.DbUsername,
		cfg.DbPassword,
		cfg.DbName,
		cfg.DbHost,
		cfg.DbPort)

	if cfg.DbPassword == "" {
		conn = fmt.Sprintf("user=%s dbname=%s host=%s port=%s sslmode=disable",
			cfg.DbUsername,
			cfg.DbName,
			cfg.DbHost,
			cfg.DbPort)
	}

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(cfg.DbMaxOpenConn)
	db.SetMaxIdleConns(cfg.DbMaxIdleConn)
	db.SetConnMaxLifetime(time.Minute * cfg.DbMaxLifeTime)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Printf("[Success] - Connected to PostgreSQL at %s:%s", cfg.DbHost, cfg.DbPort)
	return db, nil
}
