package handlers

import (
	"fmt"
	"net/http"
	"user-svc/helpers/fault"
	"user-svc/helpers/response"
	"user-svc/model"
	usecases "user-svc/usecases/user"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	user usecases.UserUsecases
}

func NewHandler(usecase usecases.UserUsecases) *Handler {
	return &Handler{
		user: usecase,
	}
}

func (h *Handler) HandleUserRegister(ctx *gin.Context) {
	var body model.RegisterUser

	if err := ctx.ShouldBindJSON(&body); err != nil {
		fault.ErrorHandler(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err),
		))
		return
	}

	bRes, err := h.user.UserRegister(body)
	if err != nil {
		fault.ErrorHandler(ctx, err)
		return
	}

	response.JSON(ctx, http.StatusAccepted, "Success", bRes)
}
