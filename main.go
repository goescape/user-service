package main

import (
	"database/sql"
	"user-svc/config"
	handlers "user-svc/handlers/user"
	repository "user-svc/repository/user"
	"user-svc/routes"
	usecases "user-svc/usecases/user"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		return
	}

	db, err := config.InitPostgreSQL(cfg.Postgres)
	if err != nil {
		return
	}
	defer db.Close()

	routes := initDepedencies(db)
	routes.SetupRoutes()
	routes.Run(cfg.Port)
}

func initDepedencies(db *sql.DB) *routes.Routes {
	userRepo := repository.NewStore(db)
	userUC := usecases.NewUserUsecase(userRepo)
	userHandler := handlers.NewHandler(userUC)

	return &routes.Routes{
		User: userHandler,
	}
}
