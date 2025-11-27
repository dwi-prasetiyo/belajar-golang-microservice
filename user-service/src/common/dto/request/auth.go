package request

type Register struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type SendOtp struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required,len=6"`
}

type VerifyRegister struct {
	Email string `json:"-" validate:"required,email"`
	Otp   string `json:"otp" validate:"required"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
