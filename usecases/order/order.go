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

	"github.com/sony/gobreaker"
)

type orderUsecase struct {
	ServiceOrderAddress string
	serverRPC           product.ProductServiceClient
	// Producer            kafkaProducer.KafkaProducerInterface
	productBreaker      *gobreaker.CircuitBreaker
	orderServiceBreaker *gobreaker.CircuitBreaker
}

func NewOrderUsecase(serviceOrderAddress string, serverRPC product.ProductServiceClient) *orderUsecase {
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
		// Producer:            p,
		productBreaker:      gobreaker.NewCircuitBreaker(cbSettings),
		orderServiceBreaker: gobreaker.NewCircuitBreaker(orderCBSettings),
	}
}

type OrderUsecases interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error)
	PaidOrder(req *model.PayOrderModel) error
	CreateOrderDenganBreaker(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error)
}

func (ou *orderUsecase) CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error) {
	var productIds string
	totalProducts := len(req.Items)

	for i, item := range req.Items {
		productIds += item.ProductId
		if i < totalProducts-1 {
			productIds += ","
		}
	}

	log.Printf("[INFO] Listing product details for IDs: %s", productIds)
	productListReq := &product.ListProductRequest{ProductIds: productIds}
	var productListResp *product.ListProductResponse

	productListResp, err := ou.serverRPC.ListProduct(ctx, productListReq)
	if err != nil {
		return nil, fmt.Errorf("failed get list product: %v", err)
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

	resp, err = client.Do(reqClient)
	if err != nil {
		log.Printf("[ERROR] Failed to send HTTP request: %v", err)
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
	return nil
}
