package model

import "time"

type Order struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	GrossAmount   int       `json:"gross_amount"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	PaymentURL    string    `json:"payment_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
