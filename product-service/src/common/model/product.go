package model

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Sku         string  `json:"sku"`
	ImageID     string  `json:"image_id"`
	Image       string  `json:"image"`
	Price       int     `json:"price"`
	Stock       int     `json:"stock"`
	Length      float32 `json:"length"`
	Width       float32 `json:"width"`
	Height      float32 `json:"height"`
	Weight      float32 `json:"weight"`
	Description string  `json:"description"`
}