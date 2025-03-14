package params

import (
	"github.com/mohammaderm/rootext/entity"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User entity.User
}

type LoginRequest struct {
	Username string
	Password string
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User   entity.User
	Tokens Tokens
}

type TokenRenewReq struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenRenewRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GetAllUserRes struct {
	Users []entity.User `json:"users"`
}
