package types

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type UserLogin struct {
	Email    string `json:"email" doc:"email" format:"email"`
	Password string `json:"password" doc:"password" minLength:"8" example:"password"`
}

type UserLoginInput struct {
	Body UserLogin
}

type UserLogoutInput struct {
	Body struct {
		ID string `json:"id" doc:"id" example:"1"`
	}
}

type LoginResponse struct {
	SetCookies []http.Cookie `header:"Set-Cookie"`
}

type Tokens struct {
	AccessToken        jwt.Token
	RefreshToken       jwt.Token
	AccessTokenString  string
	RefreshTokenStirng string
}

type key int

var UserKey key
