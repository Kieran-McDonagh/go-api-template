package endpoints

import (
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/services"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterEndpoints(api huma.API, s services.ServiceLayer) {
	RegisterUsersEndpoints(api, s.UserService)
	RegisterAuthEndpoints(api, s.UserService, s.AuthService)
}
