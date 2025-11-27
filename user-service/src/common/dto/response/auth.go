package response

import "user-service/src/common/model"

type Login struct {
	Data         *model.User
	AccessToken  string
	RefreshToken string
}

type Token struct {
	AccessToken  string
	RefreshToken string
}