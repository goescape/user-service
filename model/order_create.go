package model

type CreateOrderReq struct {
	UserId string      `json:"user_id" binding:"required"`
	Items  []OrderItem `json:"items" binding:"required"`
}

type OrderItem struct {
	ProductId string  `json:"product_id" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Qty       int64   `json:"qty" binding:"required"`
}

type CreateOrderResp struct {
	OrderId string `json:"order_id"`
}
