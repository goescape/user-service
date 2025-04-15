package handlers

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"user-svc/helpers/fault"
	"user-svc/helpers/jwt"
	"user-svc/helpers/response"
	"user-svc/model"
	usecases "user-svc/usecases/order"

	"github.com/gin-gonic/gin"
)

var (
	lockOrderRequest sync.Mutex
)

type OrderHandler struct {
	service usecases.OrderUsecases
}

func NewOrderHandler(service usecases.OrderUsecases) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	lockOrderRequest.Lock()
	defer lockOrderRequest.Unlock()

	var req model.CreateOrderReq
	authHeader := c.GetHeader("Authorization")

	// Check if the authorization header is present

	if authHeader == "" {
		log.Default().Println("Authorization header is missing")

		fault.Response(c, fault.Custom(
			http.StatusUnauthorized,
			fault.ErrUnauthorized,
			"Authorization header is missing",
		))
		return
	}

	// process the token
	token := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := jwt.GetClaims(token)
	if err != nil {
		log.Default().Println("Failed to parse token claims:", err)

		fault.Response(c, fault.Custom(
			http.StatusUnauthorized,
			fault.ErrUnauthorized,
			"Failed to parse token claims",
		))
		return
	}

	req.UserId = claims.UserId

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Default().Println("error binding JSON:", err)

		fault.Response(c, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			"failed to bind JSON: "+err.Error(),
		))
		return
	}

	resp, err := h.service.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		log.Default().Println("error creating order:", err)

		fault.Response(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, "Success", resp)
}
