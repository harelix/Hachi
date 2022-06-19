package auth

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rills-ai/Hachi/pkg/config"
	"golang.org/x/exp/slices"
	"net/http"
)

func AuthenticationMiddleware() echo.MiddlewareFunc {

	authConfig := config.New().Service.DNA.API.Auth

	// initialize JWT middleware instance
	return middleware.JWTWithConfig(middleware.JWTConfig{
		// skip public endpoints
		Skipper: func(context echo.Context) bool {
			allowedPaths := []string{
				"/health",
				"/metrics",
				"/swagger",
				"/api/v1/RELIX-2022",
			}
			return slices.Contains(allowedPaths, context.Path())
		},
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			if auth == "" {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
			}
			ks, err := keySet(authConfig.Provider)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("failed to fetch keyset: %w", err).Error())
			}

			token, err := jwt.Parse([]byte(auth), jwt.WithKeySet(ks), jwt.WithVerify(true))
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("failed to parse token: %w", err).Error())
			}

			return token, nil
		},
		ErrorHandler: func(err error) error {
			if errors.Is(err, middleware.ErrJWTMissing) {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication token")
			}
			return err
		},
	})
}
