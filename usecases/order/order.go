package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"user-svc/helpers/broker"
	"user-svc/helpers/fault"
	"user-svc/model"
	"user-svc/proto/product"
)

type orderUsecase struct {
	serviceOrderAddress string
	serverRPC           product.ProductServiceClient
	kafka               broker.KafkaProducer
}

func NewOrderUsecase(serviceOrderAddress string, serverRPC product.ProductServiceClient, kafka broker.KafkaProducer) *orderUsecase {
	return &orderUsecase{
		serviceOrderAddress: serviceOrderAddress,
		serverRPC:           serverRPC,
		kafka:               kafka,
	}
}

type OrderUsecases interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error)
	PaidOrder(req *model.PaidOrderRequest) error
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

	// Call the product service to check if the products are available
	productListReq := &product.ListProductRequest{
		ProductIds: productIds,
	}

	productListResp, err := ou.serverRPC.ListProduct(ctx, productListReq)
	if err != nil {
		log.Default().Println("Failed to call product service:", err)
		return nil, err
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

	resp, err := client.Do(reqClient)
	if err != nil {
		log.Default().Println("Failed to send request:", err)
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
