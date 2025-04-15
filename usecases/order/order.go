package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"user-svc/model"
)

type OrderUsecases interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error)
}

type orderUsecase struct {
	ServiceOrderAddress string
}

func NewOrderUsecase(serviceOrderAddress string) *orderUsecase {
	return &orderUsecase{
		ServiceOrderAddress: serviceOrderAddress,
	}
}

func (ou *orderUsecase) CreateOrder(ctx context.Context, req *model.CreateOrderReq) (*model.CreateOrderResp, error) {
	url := ou.ServiceOrderAddress + "/api/order/create"

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

	log.Default().Println("Received response:", respBody)

	return &respBody, nil
}
