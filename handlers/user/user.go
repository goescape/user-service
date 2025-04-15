package handlers

import (
	"fmt"
	"log"
	"net/http"
	"user-svc/helpers/fault"
	"user-svc/helpers/response"
	"user-svc/model"
	usecases "user-svc/usecases/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user usecases.UserUsecases
}

func NewUserHandler(usecase usecases.UserUsecases) *UserHandler {
	return &UserHandler{
		user: usecase,
	}
}

func (h *UserHandler) HandleUserRegister(ctx *gin.Context) {
	var body model.RegisterUser

	if err := ctx.ShouldBindJSON(&body); err != nil {
		fault.Response(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err),
		))
		return
	}

	bRes, err := h.user.Register(body)
	if err != nil {
		fault.Response(ctx, err)
		return
	}

	response.JSON(ctx, http.StatusAccepted, "Success", bRes)
}

func (h *UserHandler) HandleUserLogin(c *gin.Context) {
	var body model.LoginUserReq

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Default().Println("error binding JSON:", err)
		fault.Response(c, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err),
		))
		return
	}

	resp, err := h.user.Login(c.Request.Context(), &body)
	if err != nil {
		log.Default().Println("error login:", err)
		fault.Response(c, err)
		return
	}

	response.JSON(c, http.StatusOK, "Success", resp)
}
