package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"user-svc/helpers/fault"
	"user-svc/helpers/response"
	"user-svc/model"
	"user-svc/proto/product"
	usecases "user-svc/usecases/product"

	"github.com/gin-gonic/gin"
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
