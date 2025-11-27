package response

type Common struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
