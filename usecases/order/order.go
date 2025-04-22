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
	"user-svc/helpers/broker"
	"user-svc/helpers/fault"
	"user-svc/model"
	"user-svc/proto/product"

	"github.com/sony/gobreaker"
)

type orderUsecase struct {
	serviceOrderAddress string
	serverRPC           product.ProductServiceClient
	kafka               broker.KafkaProducer
	productBreaker      *gobreaker.CircuitBreaker
	orderBreaker        *gobreaker.CircuitBreaker
}

func NewOrderUsecase(serviceOrderAddress string, serverRPC product.ProductServiceClient, kafka broker.KafkaProducer, breaker gobreaker.Settings) *orderUsecase {
	orderCB := breaker
	orderCB.Name = "OrderServiceBreaker"

	return &orderUsecase{
		serviceOrderAddress: serviceOrderAddress,
		serverRPC:           serverRPC,
		kafka:               kafka,
		productBreaker:      gobreaker.NewCircuitBreaker(breaker),
		orderBreaker:        gobreaker.NewCircuitBreaker(orderCB),
	}
}

type OrderUsecases interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error)
	PaidOrder(req *model.PaidOrderRequest) error
	CreateOrderWithBreaker(ctx context.Context, body model.CreateOrderReq) (*model.CreateOrderResp, error)
}

func (ou *orderUsecase) CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error) {
	var productIds string
	totalProducts := len(req.Items)
	reduceProductQty := make([]*product.ProductItem, 0, totalProducts)

	for i, item := range req.Items {
		productIds += item.ProductId
		if i < totalProducts-1 {
			productIds += ","
		}

		// Prepare the product item for reducing quantity
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

	// Check if the products are available
	for _, item := range req.Items {
		productAvailable := false
		for _, product := range productListResp.Items {
			if item.ProductId == product.Id {
				if int64(product.Qty) < item.Qty {
					log.Default().Println("Product is out of stock:", product.Id)
					return nil, errors.New("product is out of stock")
				}
				productAvailable = true
				break
			}
		}
		if !productAvailable {
			log.Default().Println("Product not found:", item.ProductId)
			return nil, errors.New("product not found")
		}
	}

	// order request
	url := ou.serviceOrderAddress + "/api/order/create"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		log.Default().Println("Failed to marshal request body:", err)
		return nil, err
	}

	reqClient, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Default().Println("Failed to create request:", err)
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
		err = json.NewDecoder(resp.Body).Decode(&respBodyFailed)
		if err != nil {
			log.Default().Println("error when decode failed resp body:", err)
			return nil, err
		}

		log.Default().Println("Received non-200 response:", resp.Status, resp.StatusCode, respBodyFailed)
		return nil, err
	}

	var respBody model.CreateOrderResp
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		log.Default().Println("error when decode resp body:", err)
		return nil, err
	}

	// call the product service to reduce the product quantity
	reduceProductReq := &product.ReduceProductsRequest{
		Items: reduceProductQty,
	}
	_, err = ou.serverRPC.ReduceProducts(ctx, reduceProductReq)
	if err != nil {
		log.Default().Println("Failed to call product service to reduce products:", err)
		return nil, err
	}

	return &respBody, nil
}

func (ou *orderUsecase) PaidOrder(req *model.PaidOrderRequest) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return fault.Custom(
			http.StatusConflict,
			fault.ErrConflict,
			fmt.Sprintf("failed to marshal request: %v", err),
		)
	}

	if err := ou.kafka.SendMessage(model.KafkaPublish{
		Topic: "payOrder",
		Key:   "task",
		Value: payload,
	}); err != nil {
		return err
	}

	return nil
}

func (ou *orderUsecase) CreateOrderWithBreaker(ctx context.Context, body model.CreateOrderReq) (*model.CreateOrderResp, error) {
	var (
		productIds       []string
		reduceProductQty = make([]*product.ProductItem, 0, len(body.Items))
	)

	// Kumpulkan product ID dan quantity
	for _, item := range body.Items {
		productIds = append(productIds, item.ProductId)
		reduceProductQty = append(reduceProductQty, &product.ProductItem{
			ProductId: item.ProductId,
			Qty:       uint32(item.Qty),
		})
	}

	joinedProductIds := strings.Join(productIds, ",")

	// Panggil service produk dengan circuit breaker
	producList, err := ou.productBreaker.Execute(func() (interface{}, error) {
		return ou.serverRPC.ListProduct(ctx, &product.ListProductRequest{
			ProductIds: joinedProductIds,
		})
	})
	if err != nil {
		log.Printf("[CB] Product Breaker Error: %v", err)
		return nil, fault.Custom(
			http.StatusServiceUnavailable,
			fault.ErrUnavailable,
			fmt.Sprintf("service unavailable: %v", err),
		)
	}

	// Logging state circuit breaker
	switch ou.productBreaker.State() {
	case gobreaker.StateClosed:
		log.Println("[CB] Product Breaker State: CLOSED")
	case gobreaker.StateOpen:
		log.Println("[CB] Product Breaker State: OPEN")
	case gobreaker.StateHalfOpen:
		log.Println("[CB] Product Breaker State: HALF-OPEN")
	}

	products, ok := producList.(*product.ListProductResponse)
	if !ok {
		return nil, fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			"invalid product list response type",
		)
	}

	// Validasi ketersediaan produk
	for _, item := range body.Items {
		var found bool
		for _, product := range products.Items {
			if item.ProductId != product.Id {
				continue
			}
			if product.Qty < uint32(item.Qty) {
				return nil, fault.Custom(
					http.StatusUnprocessableEntity,
					fault.ErrUnprocessable,
					"product out of stock",
				)
			}
			found = true
			break
		}

		if !found {
			return nil, fault.Custom(
				http.StatusNotFound,
				fault.ErrNotFound,
				"product not found",
			)
		}
	}

	// Marshal body order
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			fmt.Sprintf("failed marshaling body: %v", err),
		)
	}

	orderURL := fmt.Sprintf("%s/api/order/create", ou.serviceOrderAddress)

	// Panggil service order dengan circuit breaker
	resOrder, err := ou.orderBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, orderURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			return nil, fault.Custom(
				http.StatusUnprocessableEntity,
				fault.ErrUnprocessable,
				fmt.Sprintf("failed to create order request: %v", err),
			)
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fault.Custom(
				http.StatusUnprocessableEntity,
				fault.ErrUnprocessable,
				fmt.Sprintf("failed to call order service: %v", err),
			)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, fault.Custom(
				http.StatusServiceUnavailable,
				fault.ErrUnavailable,
				"order service returned non-200",
			)
		}

		var apiResponse model.CreateOrderResp
		if err := json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
			return nil, fault.Custom(
				http.StatusInternalServerError,
				fault.ErrInternalServer,
				fmt.Sprintf("failed to decode order response: %v", err),
			)
		}

		return apiResponse, nil
	})
	if err != nil {
		return nil, fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			fmt.Sprintf("failed to call order service: %v", err),
		)
	}

	order := resOrder.(model.CreateOrderResp)

	// Kurangi stok produk
	_, err = ou.serverRPC.ReduceProducts(ctx, &product.ReduceProductsRequest{
		Items: reduceProductQty,
	})
	if err != nil {
		log.Printf("failed to reduce product stock: %v", err)
		// Optional: bisa return error atau lanjut tergantung business rule
	}

	return &order, nil
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
