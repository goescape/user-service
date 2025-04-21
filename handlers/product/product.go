package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"user-svc/helpers/fault"
	helperjwt "user-svc/helpers/jwt"
	"user-svc/helpers/response"
	"user-svc/model"
	"user-svc/proto/product"
	usecases "user-svc/usecases/product"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Handler struct {
	service usecases.ProductUsecases
}

func NewProductHandler(service usecases.ProductUsecases) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InsertProduct(ctx *gin.Context) {

	// Ekstrak token dari header Authorization
	tokenString, err := h.extractBearerToken(ctx)
	if err != nil {
		fault.Response(ctx, err)
		return
	}

	// Verifikasi token JWT
	dataClaims, err := helperjwt.VerifyToken(tokenString)
	if err != nil {
		// Tangani error jika token tidak valid atau kadaluarsa
		if err == jwt.ErrSignatureInvalid || strings.Contains(err.Error(), "token is expired") {
			fault.Response(ctx, err)
			return
		}
		fault.Response(ctx, err)
		return
	}
	var req model.ProductInsertReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fault.Response(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err),
		))
		return
	}

	bRes, err := h.service.InsertProduct(ctx, &product.ProductInsertRequest{
		UserId:      dataClaims.UserId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Qty:         uint32(req.Qty),
	})
	if err != nil {
		fault.Response(ctx, err)
		return
	}

	response.JSON(ctx, http.StatusCreated, "Success", bRes)
}

func (h *Handler) ListProduct(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		fault.Response(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("invalid or missing page number: %v", err.Error()),
		))
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		fault.Response(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("invalid or missing limit number: %v", err.Error()),
		))
		return
	}

	res, err := h.service.ListProduct(ctx, &product.ListProductRequest{
		Page:  uint32(page),
		Limit: uint32(limit),
	})
	if err != nil {
		fault.Response(ctx, err)
		return
	}

	response.JSON(ctx, http.StatusOK, "Success", res)
}

func (h *Handler) extractBearerToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization token")
	}

	// Pastikan format header "Authorization" adalah "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}
