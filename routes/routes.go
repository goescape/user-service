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

// Routes menyimpan semua handler dan router utama
type Routes struct {
	Router  *gin.Engine
	User    *handlers.UserHandler
	Product *productHandlers.Handler
	Order   *orderHandlers.Handler
}

// Setup menginisialisasi Gin engine, middleware, dan rute API utama
func (r *Routes) Setup(baseURL string) {
	r.Router = gin.New()

	// Tambahkan middleware CORS dan logger
	r.Router.Use(middlewares.EnabledCORS(), middlewares.Logger(r.Router))

	// Normalisasi baseURL
	if baseURL != "" && baseURL != "/" {
		baseURL = "/" + strings.Trim(baseURL, "/")
	} else {
		baseURL = "/"
	}

	// Setup semua API route
	r.setupAPIRoutes(baseURL)
}

// setupAPIRoutes mengatur route-group API utama berdasarkan baseURL
func (r *Routes) setupAPIRoutes(baseURL string) {
	apiGroup := r.Router.Group(baseURL)

	// Konfigurasi masing-masing rute
	r.configureUserRoutes(apiGroup)
	r.configureProductRoutes(apiGroup)
	r.configureOrderRoutes(apiGroup)
}

// configureUserRoutes mengatur semua endpoint terkait user
func (r *Routes) configureUserRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/user")
	userGroup.POST("/register", r.User.HandleUserRegister)
	userGroup.POST("/login", r.User.HandleUserLogin)
}

// configureProductRoutes mengatur semua endpoint terkait product
func (r *Routes) configureProductRoutes(router *gin.RouterGroup) {
	productGroup := router.Group("/product")
	productGroup.POST("/insert", r.Product.InsertProduct)
	productGroup.GET("/list", r.Product.ListProduct)
}

// configureOrderRoutes mengatur semua endpoint terkait order
func (r *Routes) configureOrderRoutes(router *gin.RouterGroup) {
	orderGroup := router.Group("/order")
	orderGroup.POST("/create", r.Order.CreateOrder)
	orderGroup.POST("/pay", r.Order.HandlePaidOrder)
}

// Run menjalankan HTTP server pada port yang ditentukan
func (r *Routes) Run(port string) {
	if r.Router == nil {
		panic("[ROUTER ERROR] Gin Engine has not been initialized. Make sure to call Setup() before Run().")
	}

	err := r.Router.Run(":" + port)
	if err != nil {
		panic(fmt.Sprintf("[SERVER ERROR] Failed to start the server on port %s: %v", port, err))
	}
}
