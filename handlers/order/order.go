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
	body.UserId = claims.UserId
	if err := ctx.ShouldBindJSON(&body); err != nil {
		fault.Response(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err.Error()),
		))
		return
	}

	bRes, err := h.service.CreateOrder(ctx.Request.Context(), &body)

	if err != nil {
		// Cek error retry khusus
		if strings.Contains(err.Error(), "service unavailable") {
			fault.Response(ctx, fault.Custom(
				http.StatusServiceUnavailable,
				fault.ErrServiceUnavailable,
				"Order service sedang tidak bisa diakses. Silakan coba beberapa saat lagi.",
			))
			return
		}

		// Error lain
		fault.Response(ctx, err)
		return
	}

	response.JSON(ctx, http.StatusCreated, "Success", bRes)
}

func (h *Handler) HandlePaidOrder(ctx *gin.Context) {
	tokenHeader := ctx.GetHeader("Authorization")
	if tokenHeader == "" || !strings.HasPrefix(tokenHeader, "Bearer ") {
		fault.Response(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			"Missing or invalid token header",
		))
		return
	}

	token := strings.TrimPrefix(tokenHeader, "Bearer ")
	_, err := jwt.GetClaims(token)
	if err != nil {
		fault.Response(ctx, fault.Custom(
			http.StatusUnauthorized,
			fault.ErrUnauthorized,
			"Failed to parse token claims",
		))
		return
	}

	var body *model.PayOrderModel
	if err := ctx.ShouldBindJSON(&body); err != nil {
		fault.Response(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err.Error()),
		))
		return
	}

	err = h.service.PaidOrder(body)
	if err != nil {
		fault.Response(ctx, err)
		return
	}

	response.JSON(ctx, http.StatusAccepted, "Success", nil)
}

func (h *Handler) CreateOrderDeanganBreaker(ctx *gin.Context) {
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
	body.UserId = claims.UserId
	if err := ctx.ShouldBindJSON(&body); err != nil {
		fault.Response(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err.Error()),
		))
		return
	}

	bRes, err := h.service.CreateOrderDenganBreaker(ctx.Request.Context(), &body)

	if err != nil {
		if strings.Contains(err.Error(), "SERVICE_UNAVAILABLE") {
			fault.Response(ctx, fault.Custom(
				http.StatusServiceUnavailable,
				fault.ErrServiceUnavailable,
				"Order service sedang tidak bisa diakses. Silakan coba beberapa saat lagi.",
			))
			return
		}

		// Error lain
		fault.Response(ctx, err)
		return
	}

	response.JSON(ctx, http.StatusCreated, "Success", bRes)
}
