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

// UserHandler merupakan handler untuk menangani permintaan terkait user
type UserHandler struct {
	user usecases.UserUsecases
}

// NewUserHandler membuat instance baru dari UserHandler
func NewUserHandler(usecase usecases.UserUsecases) *UserHandler {
	return &UserHandler{
		user: usecase,
	}
}

// HandleUserRegister menangani endpoint registrasi user
func (h *UserHandler) HandleUserRegister(ctx *gin.Context) {
	var body model.RegisterUser

	// Binding body JSON ke struct RegisterUser
	if err := ctx.ShouldBindJSON(&body); err != nil {
		fault.Response(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err),
		))
		return
	}

	// Memanggil service untuk proses registrasi user
	bRes, err := h.user.Register(body)
	if err != nil {
		fault.Response(ctx, err)
		return
	}

	// Mengembalikan response sukses
	response.JSON(ctx, http.StatusAccepted, "Success", bRes)
}

// HandleUserLogin menangani endpoint login user
func (h *UserHandler) HandleUserLogin(c *gin.Context) {
	var body model.UserLogin

	// Binding body JSON ke struct UserLogin
	if err := c.ShouldBindJSON(&body); err != nil {
		fault.Response(c, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err),
		))
		return
	}

	// Memanggil service untuk proses login user
	resp, err := h.user.Login(body)
	if err != nil {
		fault.Response(c, err)
		return
	}

	// Mengembalikan response sukses
	response.JSON(c, http.StatusOK, "Success", resp)
}
