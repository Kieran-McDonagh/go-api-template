package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/middleware"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/services"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/types"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterUsersEndpoints(api huma.API, s services.UserService) {
	huma.Register(api, huma.Operation{
		OperationID:   "post-user",
		Method:        http.MethodPost,
		Path:          "/users",
		Summary:       "Post a user",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusCreated,
		Middlewares:   huma.Middlewares{},
	}, func(ctx context.Context, i *types.NewUserInput) (*types.NewUserResponse, error) {
		user := i.Body
		id, err := s.Create(user)
		if err != nil {
			if err == types.ErrNotUniqueEmail {
				return nil, huma.Error409Conflict("a user with that email already exists")
			}
			return nil, err
		}
		resp := &types.NewUserResponse{}
		resp.Body.ID = *id
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-user-by-id",
		Method:      http.MethodGet,
		Path:        "/users/{ID}",
		Summary:     "Get a user by ID",
		Description: "Get a user by ID.",
		Tags:        []string{"Users"},
		Middlewares: huma.Middlewares{
			middleware.AuthMiddleware(api),
		},
	}, func(ctx context.Context, i *struct {
		ID string `path:"ID" example:"1" doc:"ID of user"`
	}) (*types.GetUserResponse, error) {
		user, err := s.One(i.ID)
		if err != nil {
			if err == types.ErrNotFound {
				return nil, huma.Error404NotFound(fmt.Sprintf("User with ID %s could not be found", i.ID))
			}
			return nil, err
		}
		resp := &types.GetUserResponse{
			Body: *user,
		}
		return resp, nil
	})
}
