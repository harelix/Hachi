package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

func initOauth(jwksURL string) (*jwk.Cache, error) {

	ctx := context.Background()
	ar := jwk.NewCache(ctx)
	err := ar.Register(jwksURL, jwk.WithMinRefreshInterval(15*time.Minute))

	if err != nil {
		return nil, err
	}

	_, err = ar.Refresh(ctx, jwksURL)
	if err != nil {
		fmt.Printf("failed to refresh JWKS: %s\n", err)
		return nil, err
	}
	return ar, nil
}

var cache *jwk.Cache

func keySet(jwksURL string) (jwk.Set, error) {
	var err error

	provider := fmt.Sprintf("%s,%s", jwksURL, "/.well-known/jwks.json")
	if cache == nil {
		cache, err = initOauth(provider)
		if err != nil {
			return nil, err
		}
	}
	return cache.Get(context.Background(), provider)
}
