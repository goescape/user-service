package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgreSQLConfig struct {
	DbHost     string
	DbPort     string
	DbUsername string
	DbPassword string
	DbName     string
}

func InitPostgreSQL(cfg PostgreSQLConfig) (*sql.DB, error) {
	conn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.DbUsername,
		cfg.DbPassword,
		cfg.DbName,
		cfg.DbHost,
		cfg.DbPort)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, nil
}
