package routes

import (
	"fmt"
	"strings"
	handlers "user-svc/handlers/user"
	"user-svc/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Routes struct {
	Router *gin.Engine
	User   *handlers.Handler
}

func (r *Routes) SetupRoutes() {
	r.Router = gin.New()
	r.Router.Use(middlewares.EnabledCORS(), middlewares.Logger(r.Router))

	r.setupAPIRoutes()
}

func (r *Routes) setupAPIRoutes() {
	baseURL := viper.GetString("BASE_URL_PATH")
	if baseURL == "" || baseURL == "/" {
		baseURL = "/"
	} else {
		baseURL = "/" + strings.TrimPrefix(baseURL, "/")
	}

	apiGroup := r.Router.Group(baseURL)
	r.configureUserRoutes(apiGroup)
}

func (r *Routes) configureUserRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/user")
	userGroup.POST("/register", r.User.HandleUserRegister)
}

func (r *Routes) Run(port string) {
	if r.Router == nil {
		panic("[ROUTER ERROR] Gin Engine has not been initialized. Make sure to call SetupRouter() before Run().")
	}

	addr := fmt.Sprintf(":%s", port)
	if err := r.Router.Run(addr); err != nil {
		panic(fmt.Sprintf("[SERVER ERROR] Failed to start the server on port %s: %v", port, err))
	}
}
