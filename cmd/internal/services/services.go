package services

import "github.com/Kieran-McDonagh/go-api-template/cmd/internal/providers"

type ServiceLayer struct {
	UserService
	AuthService
}

func NewServiceLayer(p providers.ProviderLayer) ServiceLayer {
	return ServiceLayer{
		UserService: NewUserService(p),
		AuthService: NewAuthService(p),
	}
}
