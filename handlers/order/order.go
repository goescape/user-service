package handlers

import (
	"fmt"
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

var mu sync.Mutex

type Handler struct {
	service usecases.OrderUsecases
}

func NewOrderHandler(service usecases.OrderUsecases) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateOrder(ctx *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	tokenHeader := ctx.GetHeader("Authorization")
	if tokenHeader == "" {
		fault.Response(ctx, fault.Custom(
			http.StatusUnauthorized,
			fault.ErrUnauthorized,
			"Authorization header is missing",
		))
		return
	}

	token := strings.TrimPrefix(tokenHeader, "Bearer ")
	claims, err := jwt.GetClaims(token)
	if err != nil {
		fault.Response(ctx, fault.Custom(
			http.StatusUnauthorized,
			fault.ErrUnauthorized,
			"Failed to parse token claims",
		))
		return
	}

	var body model.CreateOrderReq
	if err := ctx.ShouldBindJSON(&body); err != nil {
		fault.Response(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err.Error()),
		))
		return
	}

	body.UserId = claims.UserId

	bRes, err := h.service.CreateOrder(ctx.Request.Context(), &body)
	if err != nil {
		fault.Response(ctx, err)
		return
	}

	response.JSON(ctx, http.StatusCreated, "Success", bRes)
}
