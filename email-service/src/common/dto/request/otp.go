package request

type SendOtp struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required"`
	Year  int    `validate:"omitempty"`
}