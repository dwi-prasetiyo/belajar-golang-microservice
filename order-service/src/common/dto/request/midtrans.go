package request

type MidtransTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int    `json:"gross_amount"`
}

type MidtransTransactionMetadata struct {
	OriginalOrderID string `json:"original_order_id"`
}

type MidtransTransaction struct {
	TransactionDetails MidtransTransactionDetails  `json:"transaction_details"`
	Metadata           MidtransTransactionMetadata `json:"metadata"`
}
