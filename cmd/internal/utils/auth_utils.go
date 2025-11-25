package utils

import (
	"log"
	"os"
	"time"

	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func GetSecretKey() []byte {
	if err := godotenv.Load("internal/.env"); err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}
	key := os.Getenv("SECRET_KEY")
	byteArray := []byte(key)
	return byteArray
}

func CreateTokens(user types.UserClaims) (*types.Tokens, error) {
	secretKey := GetSecretKey()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":      uuid.New(),
			"user_id": user.ID,
			"email":   user.Email,
			"role":    user.Role,
			"exp":     time.Now().Add(time.Minute * 15).Unix(),
		})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":      uuid.New(),
			"user_id": user.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	tokens := types.Tokens{
		AccessToken:        *accessToken,
		RefreshToken:       *refreshToken,
		AccessTokenString:  accessTokenString,
		RefreshTokenStirng: refreshTokenString,
	}

	return &tokens, nil
}
