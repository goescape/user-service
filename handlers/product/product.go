package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"user-svc/helpers/fault"
	"user-svc/helpers/response"
	"user-svc/model"
	"user-svc/proto/product"
	usecases "user-svc/usecases/product"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service usecases.ProductUsecases
}

func NewProductHandler(service usecases.ProductUsecases) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) InsertProduct(c *gin.Context) {
	var req model.ProductInsertReq

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Default().Println("error binding JSON:", err)
		fault.Response(c, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err),
		))
		return
	}

	rpcReq := &product.ProductInsertRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Qty:         uint32(req.Qty),
	}

	res, err := h.service.InsertProduct(c.Request.Context(), rpcReq)
	if err != nil {
		log.Default().Println("error inserting product:", err)
		fault.Response(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, "Success", res)
}

func (h *ProductHandler) ListProduct(c *gin.Context) {
	// get query param for page and limit
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// convert page and limit to int
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Default().Println("error converting page to int:", err)
		fault.Response(c, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("invalid page number: %v", err),
		))
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Default().Println("error converting limit to int:", err)
		fault.Response(c, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("invalid limit number: %v", err),
		))
		return
	}

	rpcReq := &product.ListProductRequest{
		Page:  uint32(pageInt),
		Limit: uint32(limitInt),
	}

	res, err := h.service.ListProduct(c.Request.Context(), rpcReq)
	if err != nil {
		log.Default().Println("error listing product:", err)
		fault.Response(c, err)
		return
	}

	response.JSON(c, http.StatusOK, "Success", res)
}
