package request

type MidtransTx struct {
	OrderID           string             `json:"order_id"`
	StatusCode        string             `json:"status_code"`
	GrossAmount       string             `json:"gross_amount"`
	TransactionStatus string             `json:"transaction_status"`
	SignatureKey      string             `json:"signature_key"`
	PaymentType       string             `json:"payment_type"`
	FraudStatus       string             `json:"fraud_status"`
	Metadata          MidtransTxMetadata `json:"metadata"`
}

type MidtransTxMetadata struct {
	OriginalOrderID string `json:"original_order_id"`
}
