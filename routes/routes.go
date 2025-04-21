package routes

import (
	"fmt"
	"strings"
	orderHandlers "user-svc/handlers/order"
	productHandlers "user-svc/handlers/product"
	handlers "user-svc/handlers/user"
	"user-svc/middlewares"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	Router  *gin.Engine
	User    *handlers.UserHandler
	Product *productHandlers.Handler
	Order   *orderHandlers.Handler
}

func (r *Routes) Setup(baseURL string) {
	r.Router = gin.New()
	r.Router.Use(middlewares.EnabledCORS(), middlewares.Logger(r.Router))

	if baseURL != "" && baseURL != "/" {
		baseURL = "/" + strings.Trim(baseURL, "/")
	} else {
		baseURL = "/"
	}

	r.setupAPIRoutes(baseURL)
}

func (r *Routes) setupAPIRoutes(baseURL string) {
	apiGroup := r.Router.Group(baseURL)
	r.configureUserRoutes(apiGroup)
	r.configureProductRoutes(apiGroup)
	r.configureOrderRoutes(apiGroup)
}

func (r *Routes) configureUserRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/user")
	userGroup.POST("/register", r.User.HandleUserRegister)
	userGroup.POST("/login", r.User.HandleUserLogin)
}

func (r *Routes) configureProductRoutes(router *gin.RouterGroup) {
	productGroup := router.Group("/product")
	productGroup.POST("/insert", r.Product.InsertProduct)
	productGroup.GET("/list", r.Product.ListProduct)
}

func (r *Routes) configureOrderRoutes(router *gin.RouterGroup) {
	orderGroup := router.Group("/order")
	orderGroup.POST("/create", r.Order.CreateOrder)
	orderGroup.POST("/pay", r.Order.HandlePaidOrder)
}

func (r *Routes) Run(port string) {
	if r.Router == nil {
		panic("[ROUTER ERROR] Gin Engine has not been initialized. Make sure to call Setup() before Run().")
	}

	err := r.Router.Run(":" + port)
	if err != nil {
		panic(fmt.Sprintf("[SERVER ERROR] Failed to start the server on port %s: %v", port, err))
	}
}
