package middleware

import (
	"log"

	"github.com/danielgtaylor/huma/v2"
)

func LoggingMiddleware(ctx huma.Context, next func(huma.Context)) {
	log.Println(ctx.Method(), ctx.Operation().Path)

	// Call the next middleware in the chain. This eventually calls the
	// operation handler as well.
	next(ctx)
}
