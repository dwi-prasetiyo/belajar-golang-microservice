package response

type MidtransTx struct {
	OrderID    string `json:"order_id"`
	PaymentURL string `json:"payment_url"`
}