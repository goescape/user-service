package main

import (
	"database/sql"
	"user-svc/config"
	orderHandlers "user-svc/handlers/order"
	productHandlers "user-svc/handlers/product"
	handlers "user-svc/handlers/user"
	"user-svc/helpers/broker"
	"user-svc/proto/product"
	repository "user-svc/repository/user"
	"user-svc/routes"
	orderUC "user-svc/usecases/order"
	productUC "user-svc/usecases/product"
	usecases "user-svc/usecases/user"

	"github.com/IBM/sarama"
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

	kafka, err := config.InitKafkaProducer(cfg.Kafka)
	if err != nil {
		return
	}
	defer (*kafka).Close()

	routes := initDepedencies(cfg, db, rpc, redis, kafka)
	routes.Setup(cfg.BaseURL)
	routes.Run(cfg.Port)
}

func initDepedencies(cfg *config.Config, db *sql.DB, rpc *grpc.ClientConn, redis *redis.Client, kafka *sarama.SyncProducer) *routes.Routes {
	userRepo := repository.NewUserStore(db)
	userUC := usecases.NewUserUsecase(userRepo, redis)
	userHandler := handlers.NewUserHandler(userUC)

	productRPC := product.NewProductServiceClient(rpc)
	productUC := productUC.NewProductUsecase(productRPC)
	productHandler := productHandlers.NewProductHandler(productUC)

	producer := broker.NewProducer(*kafka)
	orderUC := orderUC.NewOrderUsecase(cfg.ServiceOrderAdress, productRPC, producer)
	orderHandler := orderHandlers.NewOrderHandler(orderUC)

	return &routes.Routes{
		User:    userHandler,
		Product: productHandler,
		Order:   orderHandler,
	}
}
