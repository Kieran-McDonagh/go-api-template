package providers

import "database/sql"

type ProviderLayer struct {
	UserProvider
}

func NewProviderLayer(DB *sql.DB) ProviderLayer {
	return ProviderLayer{
		UserProvider: NewUserProvider(DB),
	}
}
