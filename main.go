package main

import (
	"database/sql"
	"user-svc/config"
	handlers "user-svc/handlers/user"
	repository "user-svc/repository/user"
	"user-svc/routes"
	usecases "user-svc/usecases/user"

	"google.golang.org/grpc"
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

	rpc, err := config.RPCDial(cfg.Grpc)
	if err != nil {
		return
	}

	routes := initDepedencies(db, rpc)
	routes.SetupRoutes()
	routes.Run(cfg.Port)
}

func initDepedencies(db *sql.DB, rpc *grpc.ClientConn) *routes.Routes {
	userRepo := repository.NewStore(db)
	userUC := usecases.NewUserUsecase(userRepo)
	userHandler := handlers.NewHandler(userUC)

	return &routes.Routes{
		User: userHandler,
	}
}
