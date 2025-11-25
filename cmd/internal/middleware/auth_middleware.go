package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/types"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/utils"
	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RoleToString(r types.UserRole) string {
	return string(r)
}

func refreshTokens(refreshTokenString string, userFromClaims types.UserClaims) (*types.Tokens, error) {
	secretKey := utils.GetSecretKey()
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !refreshToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	newTokens, err := utils.CreateTokens(userFromClaims)
	if err != nil {
		return nil, err
	}

	tokens := types.Tokens{
		AccessToken:        newTokens.AccessToken,
		RefreshToken:       newTokens.RefreshToken,
		AccessTokenString:  newTokens.AccessTokenString,
		RefreshTokenStirng: newTokens.RefreshTokenStirng,
	}

	return &tokens, nil
}

func verifyToken(accessTokenString, refreshTokenString string, ctx huma.Context) (*types.UserClaims, error) {
	secretKey := utils.GetSecretKey()

	// Parse without validating so we can read claims on expired tokens
	accessToken, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithoutClaimsValidation())

	if accessToken == nil || accessToken.Claims == nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims := accessToken.Claims.(jwt.MapClaims)

	exp, _ := accessToken.Claims.GetExpirationTime()
	expired := exp.Time.Before(time.Now())

	if !expired {
		// Token valid → return user
		user := types.UserClaims{
			ID:    claims["user_id"].(string),
			Email: claims["email"].(string),
			Role:  claims["role"].(string),
		}
		return &user, nil
	}

	// Token expired → refresh
	user := types.UserClaims{
		ID:    claims["user_id"].(string),
		Email: claims["email"].(string),
		Role:  claims["role"].(string),
	}

	newTokens, err := refreshTokens(refreshTokenString, user)
	if err != nil {
		return nil, err
	}

	// Set cookies
	ctx.AppendHeader("Set-Cookie", (&http.Cookie{
		Name:   "access",
		Value:  newTokens.AccessTokenString,
		Secure: true,
		Path:   "/",
	}).String())

	ctx.AppendHeader("Set-Cookie", (&http.Cookie{
		Name:   "refresh",
		Value:  newTokens.RefreshTokenStirng,
		Secure: true,
		Path:   "/",
	}).String())

	return &user, nil
}

func AuthMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		accessToken, err := huma.ReadCookie(ctx, "access")
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthorised")
			return
		}

		refreshToken, err := huma.ReadCookie(ctx, "refresh")
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthorised")
			return
		}

		userFromClaims, err := verifyToken(accessToken.Value, refreshToken.Value, ctx)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthorised")
			return
		}

		base := ctx.Context()
		base = context.WithValue(base, utils.UserClaimsKey, userFromClaims)
		ctx = huma.WithContext(ctx, base)
		next(ctx)
	}
}
