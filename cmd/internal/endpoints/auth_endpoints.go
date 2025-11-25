package endpoints

import (
	"context"
	"net/http"

	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/middleware"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/services"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/types"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/utils"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterAuthEndpoints(api huma.API, u services.UserService, a services.AuthService) {
	huma.Register(api, huma.Operation{
		OperationID:   "login",
		Method:        http.MethodPost,
		Path:          "/login",
		Summary:       "Login",
		Tags:          []string{"Authentication"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *types.UserLoginInput) (*types.LoginResponse, error) {
		tokens, err := a.Login(i.Body.Email, i.Body.Password)
		if err != nil {
			if err == types.ErrIncorrectLogin {
				return nil, huma.Error401Unauthorized("Incorrect email or password")
			}
			return nil, huma.Error401Unauthorized("Error authenticating user")
		}

		resp := &types.LoginResponse{
			SetCookies: []http.Cookie{
				{
					Name:   "access",
					Value:  tokens.AccessTokenString,
					Secure: true,
					Path:   "/",
				},
				{
					Name:   "refresh",
					Value:  tokens.RefreshTokenStirng,
					Secure: true,
					Path:   "/",
				},
			},
		}
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "logout",
		Method:        http.MethodPost,
		Path:          "/logout",
		Summary:       "Logout",
		Tags:          []string{"Authentication"},
		DefaultStatus: http.StatusOK,
		Middlewares: huma.Middlewares{
			middleware.AuthMiddleware(api),
		},
	}, func(ctx context.Context, i *types.UserLogoutInput) (*types.LoginResponse, error) {
		userFromContext, _ := utils.UserClaimsFromContext(ctx)

		if userFromContext.ID == i.Body.ID {
			resp := &types.LoginResponse{
				SetCookies: []http.Cookie{
					{
						Name:   "access",
						Value:  "",
						Secure: true,
						Path:   "/",
						MaxAge: -1,
					},
					{
						Name:   "refresh",
						Value:  "",
						Secure: true,
						Path:   "/",
						MaxAge: -1,
					},
				},
			}
			return resp, nil
		}

		return nil, huma.Error401Unauthorized("Error logging out user")
	})
}
