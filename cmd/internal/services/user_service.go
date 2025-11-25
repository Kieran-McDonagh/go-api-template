package services

import (
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/providers"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	providers.ProviderLayer
}

func NewUserService(p providers.ProviderLayer) UserService {
	return UserService{
		ProviderLayer: p,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (U UserService) Create(newUser types.NewUser) (*string, error) {
	hashedPassword, err := hashPassword(newUser.Password)
	if err != nil {
		return nil, err
	}

	newUser.Password = hashedPassword
	id, err := U.UserProvider.Create(newUser)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (U UserService) One(ID string) (*types.GetUserResponseBody, error) {
	result, err := U.UserProvider.One(ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
