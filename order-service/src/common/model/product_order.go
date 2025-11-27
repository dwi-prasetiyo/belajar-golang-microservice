package model

type ProductOrder struct {
	ID        int    `json:"id"`
	OrderID   string `json:"order_id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
}
