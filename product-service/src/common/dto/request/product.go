package request

type CreateProduct struct {
	Name        string  `json:"name" validate:"required"`
	Sku         string  `json:"sku" validate:"required"`
	ImageID     string  `json:"-" validate:"required"`
	Image       string  `json:"-" validate:"required"`
	Price       int     `json:"price" validate:"required"`
	Stock       int     `json:"stock" validate:"required"`
	Length      float32 `json:"length" validate:"required"`
	Width       float32 `json:"width" validate:"required"`
	Height      float32 `json:"height" validate:"required"`
	Weight      float32 `json:"weight" validate:"required"`
	Description string  `json:"description" validate:"required"`
}

type FindManyProduct struct {
	Page   int    `json:"page" validate:"min=1"`
	Search string `json:"search"`
}