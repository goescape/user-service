package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"user-svc/model"
	"user-svc/proto/product"

	kafkaProducer "user-svc/broker/kafka/producer"

	"github.com/sony/gobreaker"
)

type orderUsecase struct {
	ServiceOrderAddress string
	serverRPC           product.ProductServiceClient
	Producer            kafkaProducer.KafkaProducerInterface
	productBreaker      *gobreaker.CircuitBreaker
	orderServiceBreaker *gobreaker.CircuitBreaker
}

func NewOrderUsecase(serviceOrderAddress string, serverRPC product.ProductServiceClient, p kafkaProducer.KafkaProducerInterface) *orderUsecase {
	cbSettings := gobreaker.Settings{
		Name:        "ProductServiceBreaker",
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3 // circuit akan open jika gagal 3 kali berturut-turut
		},
	}

	orderCBSettings := cbSettings
	orderCBSettings.Name = "OrderServiceBreaker" // breaker berbeda untuk service order dan product

	return &orderUsecase{
		ServiceOrderAddress: serviceOrderAddress,
		serverRPC:           serverRPC,
		Producer:            p,
		productBreaker:      gobreaker.NewCircuitBreaker(cbSettings),
		orderServiceBreaker: gobreaker.NewCircuitBreaker(orderCBSettings),
	}
}

type OrderUsecases interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error)
	PaidOrder(req *model.PayOrderModel) error
	CreateOrderDenganBreaker(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error)
}

// dengan retry
func (ou *orderUsecase) CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error) {
	var productIds string
	totalProducts := len(req.Items)
	reduceProductQty := make([]*product.ProductItem, 0, totalProducts)

	for i, item := range req.Items {
		productIds += item.ProductId
		if i < totalProducts-1 {
			productIds += ","
		}
		reduceProductQty = append(reduceProductQty, &product.ProductItem{
			ProductId: item.ProductId,
			Qty:       uint32(item.Qty),
		})
	}

	log.Printf("[INFO] Listing product details for IDs: %s", productIds)
	productListReq := &product.ListProductRequest{ProductIds: productIds}
	var productListResp *product.ListProductResponse

	err := retry(2, 40*time.Second, func() error {
		var callErr error
		productListResp, callErr = ou.serverRPC.ListProduct(ctx, productListReq)
		return callErr
	})
	if err != nil {
		log.Printf("[ERROR] Failed to call ListProduct after retries: %v", err)
		return nil, errors.New("SERVICE_UNAVAILABLE")
	}

	for _, item := range req.Items {
		productAvailable := false
		for _, product := range productListResp.Items {
			if item.ProductId == product.Id {
				if int64(product.Qty) < item.Qty {
					log.Printf("[ERROR] Product %s out of stock. Requested: %d, Available: %d", product.Id, item.Qty, product.Qty)
					return nil, errors.New("product is out of stock")
				}
				productAvailable = true
				break
			}
		}
		if !productAvailable {
			log.Printf("[ERROR] Product not found: %s", item.ProductId)
			return nil, errors.New("product not found")
		}
	}

	url := ou.ServiceOrderAddress + "/api/order/create"
	log.Printf("[INFO] Sending create order request to %s", url)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal request body: %v", err)
		return nil, err
	}

	reqClient, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("[ERROR] Failed to create HTTP request: %v", err)
		return nil, err
	}
	reqClient = reqClient.WithContext(ctx)
	reqClient.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	var resp *http.Response

	err = retry(1, 25*time.Second, func() error {
		resp, err = client.Do(reqClient)
		return err
	})
	if err != nil {
		log.Printf("[ERROR] Failed to send HTTP request after retries: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respBodyFailed any
		if err := json.NewDecoder(resp.Body).Decode(&respBodyFailed); err != nil {
			log.Printf("[ERROR] Failed to decode error response: %v", err)
			return nil, err
		}
		log.Printf("[ERROR] Non-200 response from order service: %d - %v", resp.StatusCode, respBodyFailed)
		return nil, fmt.Errorf("non-200 response from order service")
	}

	var respBody model.CreateOrderResp
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		log.Printf("[ERROR] Failed to decode success response: %v", err)
		return nil, err
	}

	log.Printf("[INFO] Order created successfully with ID: %s", respBody.OrderId)

	reduceProductReq := &product.ReduceProductsRequest{Items: reduceProductQty}
	err = retry(2, 25*time.Second, func() error {
		_, callErr := ou.serverRPC.ReduceProducts(ctx, reduceProductReq)
		return callErr
	})
	if err != nil {
		log.Printf("[ERROR] Failed to reduce product quantity after retries: %v", err)
		return nil, err
	}

	log.Printf("[INFO] Product quantities reduced successfully")
	return &respBody, nil
}

// dengan breaker
func (ou *orderUsecase) CreateOrderDenganBreaker(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error) {
	// Persiapan ID produk dan payload reduce stock
	var productIds []string
	reduceProductQty := make([]*product.ProductItem, 0, len(req.Items))

	for _, item := range req.Items {
		productIds = append(productIds, item.ProductId)
		reduceProductQty = append(reduceProductQty, &product.ProductItem{
			ProductId: item.ProductId,
			Qty:       uint32(item.Qty),
		})
	}

	// --- Circuit Breaker untuk Product Service ---
	productListReq := &product.ListProductRequest{
		ProductIds: strings.Join(productIds, ","),
	}

	productListRespIface, err := ou.productBreaker.Execute(func() (interface{}, error) {
		return ou.serverRPC.ListProduct(ctx, productListReq) // gRPC call dibungkus circuit breaker
	})

	state := ou.productBreaker.State()
	switch state {
	case gobreaker.StateClosed:
		log.Println("[CB] Product Breaker State: CLOSED")
	case gobreaker.StateOpen:
		log.Println("[CB] Product Breaker State: OPEN")
	case gobreaker.StateHalfOpen:
		log.Println("[CB] Product Breaker State: HALF-OPEN")
	}
	if err != nil {
		log.Println("[CB] Failed to call product service:", err)
		return nil, errors.New("SERVICE_UNAVAILABLE")
	}
	productListResp := productListRespIface.(*product.ListProductResponse)

	// Validasi apakah semua produk tersedia dan mencukupi jumlahnya
	for _, item := range req.Items {
		found := false
		for _, p := range productListResp.Items {
			if item.ProductId == p.Id {
				if int64(p.Qty) < item.Qty {
					log.Println("Product out of stock:", p.Id)
					return nil, errors.New("product out of stock")
				}
				found = true
				break
			}
		}
		if !found {
			log.Println("Product not found:", item.ProductId)
			return nil, errors.New("product not found")
		}
	}

	// --- Kirim HTTP Request ke Order Service ---
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		log.Println("Failed to marshal order request:", err)
		return nil, err
	}

	orderURL := ou.ServiceOrderAddress + "/api/order/create"
	orderRespIface, err := ou.orderServiceBreaker.Execute(func() (interface{}, error) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, orderURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(httpReq) // HTTP call dibungkus circuit breaker
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var failResp any
			json.NewDecoder(resp.Body).Decode(&failResp)
			log.Println("Order service returned non-200:", resp.StatusCode, failResp)
			return nil, errors.New("SERVICE_UNAVAILABLE")
		}

		var successResp model.CreateOrderResp
		err = json.NewDecoder(resp.Body).Decode(&successResp)
		if err != nil {
			return nil, err
		}
		return &successResp, nil
	})
	if err != nil {
		log.Println("[CB] Failed to call order service:", err)
		return nil, err
	}
	orderResp := orderRespIface.(*model.CreateOrderResp)

	// --- Reduce Stock (tanpa circuit breaker, bisa dipertimbangkan ditambah) ---
	reduceReq := &product.ReduceProductsRequest{Items: reduceProductQty}
	if _, err := ou.serverRPC.ReduceProducts(ctx, reduceReq); err != nil {
		log.Println("Failed to reduce product stock:", err)
		// Tidak return error agar order tetap sukses
	}

	return orderResp, nil
}

func (ou *orderUsecase) PaidOrder(req *model.PayOrderModel) error {
	newData, _ := json.Marshal(req)                             // Serialize request ke JSON
	err := ou.Producer.SendMessage("payOrder", "task", newData) // Kirim ke Kafka
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func retry(attempts int, sleep time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil // sukses
		}
		log.Printf("[Retry %d/%d] Error: %v", i+1, attempts, err)
		time.Sleep(sleep)
		sleep *= 2 // exponential backoff
	}
	return fmt.Errorf("after %d attempts, last error: %w", attempts, err)
}
