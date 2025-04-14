package main

import (
	"database/sql"
	"user-svc/config"
	productHandlers "user-svc/handlers/product"
	handlers "user-svc/handlers/user"
	"user-svc/proto/product"
	repository "user-svc/repository/user"
	"user-svc/routes"
	productUsecases "user-svc/usecases/product"
	usecases "user-svc/usecases/user"

	"github.com/redis/go-redis/v9"
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

	redis, err := config.InitRedis(cfg.Redis)
	if err != nil {
		return
	}
	defer redis.Close()

	rpc, err := config.RPCDial(cfg.Grpc)
	if err != nil {
		return
	}

	routes := initDepedencies(db, rpc, redis)
	routes.Setup(cfg.BaseURL)
	routes.Run(cfg.Port)
}

func initDepedencies(db *sql.DB, rpc *grpc.ClientConn, redis *redis.Client) *routes.Routes {
	userRepo := repository.NewUserStore(db)
	userUC := usecases.NewUserUsecase(userRepo, redis)
	userHandler := handlers.NewUserHandler(userUC)

	productGRPCClient := product.NewProductServiceClient(rpc)
	productUC := productUsecases.NewProductUsecase(productGRPCClient)
	productHandler := productHandlers.NewProductHandler(productUC)

	return &routes.Routes{
		User:    userHandler,
		Product: productHandler,
	}
}
