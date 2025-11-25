package internal

import (
	"fmt"
	"net/http"

	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/database"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/endpoints"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/middleware"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/providers"
	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/services"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

func Start() {
	db := database.InitDB()
	defer db.Close()

	p := providers.NewProviderLayer(db)
	s := services.NewServiceLayer(p)
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))
		api.UseMiddleware(middleware.LoggingMiddleware)
		endpoints.RegisterEndpoints(api, s)
		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			fmt.Printf("Open API Docs available at http://localhost:%d/docs#\n", options.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// TODO: Store a blacklist of invalid refresh tokens in an an in-memory cache to use when validating refresh tokens. e.g. when a user resets their password blacklist their refresh token.
	cli.Run()
}
