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
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

func main() {
	// Load konfigurasi dari file/env
	cfg, err := config.Load()
	if err != nil {
		return
	}

	// Inisialisasi koneksi ke PostgreSQL
	db, err := config.InitPostgreSQL(cfg.Postgres)
	if err != nil {
		return
	}
	defer db.Close()

	// Inisialisasi koneksi ke Redis
	redis, err := config.InitRedis(cfg.Redis)
	if err != nil {
		return
	}
	defer redis.Close()

	// Inisialisasi koneksi ke gRPC service
	rpc, err := config.RPCDial(cfg.Grpc)
	if err != nil {
		return
	}

	// Inisialisasi Kafka producer
	kafka, err := config.InitKafkaProducer(cfg.Kafka)
	if err != nil {
		return
	}
	defer (*kafka).Close()

	// Inisialisasi circuit breaker
	breaker := config.InitBreaker()

	// Inisialisasi seluruh dependency dan routing
	routes := initDepedencies(cfg, db, rpc, redis, kafka, breaker)
	routes.Setup(cfg.BaseURL) // Setup routing base URL
	routes.Run(cfg.Port)      // Jalankan HTTP server di port yang ditentukan
}

// initDepedencies menginisialisasi seluruh dependency service dan mengembalikan instance Routes
func initDepedencies(
	cfg *config.Config,
	db *sql.DB,
	rpc *grpc.ClientConn,
	redis *redis.Client,
	kafka *sarama.SyncProducer,
	breaker *gobreaker.Settings,
) *routes.Routes {

	// User service
	userRepo := repository.NewUserStore(db)
	userUC := usecases.NewUserUsecase(userRepo, redis)
	userHandler := handlers.NewUserHandler(userUC)

	// Product service (gRPC client)
	productRPC := product.NewProductServiceClient(rpc)
	productUC := productUC.NewProductUsecase(productRPC)
	productHandler := productHandlers.NewProductHandler(productUC)

	// Order service dengan Kafka producer
	producer := broker.NewProducer(*kafka)
	orderUC := orderUC.NewOrderUsecase(cfg.ServiceOrderAdress, productRPC, producer, *breaker)
	orderHandler := orderHandlers.NewOrderHandler(orderUC)

	// Kembalikan semua handler yang sudah diinject ke struct Routes
	return &routes.Routes{
		User:    userHandler,
		Product: productHandler,
		Order:   orderHandler,
	}
}
