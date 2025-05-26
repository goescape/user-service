package main

import (
	"database/sql"
	"log"
	"user-svc/config"
	orderHandlers "user-svc/handlers/order"
	productHandlers "user-svc/handlers/product"
	handlers "user-svc/handlers/user"
	"user-svc/proto/product"
	repository "user-svc/repository/user"
	"user-svc/routes"
	orderUC "user-svc/usecases/order"
	productUC "user-svc/usecases/product"
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

	log.Println("cek")
	rpc, err := config.RPCDial(cfg.Grpc)
	if err != nil {
		return
	}
	log.Println("cek2")

	routes := initDepedencies(cfg, db, rpc, redis)
	routes.Setup(cfg.BaseURL)
	routes.Run(cfg.Port)
}

func initDepedencies(cfg *config.Config, db *sql.DB, rpc *grpc.ClientConn, redis *redis.Client) *routes.Routes {
	userRepo := repository.NewUserStore(db)
	userUC := usecases.NewUserUsecase(userRepo, redis)
	userHandler := handlers.NewUserHandler(userUC)

	productRPC := product.NewProductServiceClient(rpc)
	productUC := productUC.NewProductUsecase(productRPC)
	productHandler := productHandlers.NewProductHandler(productUC)

	orderUC := orderUC.NewOrderUsecase(cfg.ServiceOrderAdress, productRPC)
	orderHandler := orderHandlers.NewOrderHandler(orderUC)

	return &routes.Routes{
		User:    userHandler,
		Product: productHandler,
		Order:   orderHandler,
	}
}
