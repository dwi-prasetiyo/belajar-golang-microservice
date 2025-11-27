package request

type RabbitMQMessage struct {
	RequestID any `json:"request_id"`
	UserID    any `json:"user_id"`
	Message   any `json:"message"`
}
