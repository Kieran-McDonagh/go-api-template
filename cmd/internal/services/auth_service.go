package services

import (
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/providers"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/types"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	providers.ProviderLayer
}

func verifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewAuthService(p providers.ProviderLayer) AuthService {
	return AuthService{
		ProviderLayer: p,
	}
}

func (A AuthService) Login(email, password string) (*types.Tokens, error) {
	user, err := A.UserProvider.OneByEmail(email)
	if err != nil {
		return nil, err
	}

	if !verifyPassword(user.Password, password) {
		return nil, types.ErrIncorrectLogin
	}

	userClaims := types.UserClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	tokens, err := utils.CreateTokens(userClaims)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
