package utils

import (
	"context"

	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/types"
)

type key int

var UserClaimsKey key = 1

func UserClaimsFromContext(ctx context.Context) (*types.UserClaims, bool) {
	u, ok := ctx.Value(UserClaimsKey).(*types.UserClaims)
	return u, ok
}
