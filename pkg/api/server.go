package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	_ "github.com/rills-ai/Hachi/pkg/api/docs"
	"github.com/rills-ai/Hachi/pkg/api/routes"
	"github.com/rills-ai/Hachi/pkg/config"
)

//StartAPIServer will panic if it cannot start the server
func StartAPIServer(ctx context.Context) {

	serverConfig := config.New().Service.DNA.API

	e := echo.New()
	//e.Use(tracing.Middleware(tracing.WithServiceName("Hachi-" + config.New().Service.Type.String())))
	e.Use(echoMw.Logger())
	// the Gzip middleware has a high cost, maybe its better applied at the routing level to a specific endpoint?
	e.Use(echoMw.GzipWithConfig(echoMw.GzipConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}
			return false
		},
	}))
	e.Use(echoMw.Recover())
	// the requestID middleware only sets the ID on the response, wont give you an ID to send to NATS
	e.Use(echoMw.RequestID())
	e.Use(echoMw.CORS())

	e.HideBanner = true
	p := prometheus.NewPrometheus("hachi", nil)
	p.Use(e)

	if serverConfig.Auth.Enabled {
		//e.Use(auth.AuthenticationMiddleware())
	}
	//routes bindings
	err := routes.RegisterRoutes(e)
	if err != nil {
		e.Logger.Panic(err)
	}

	quit := make(chan os.Signal, 1)

	//todo: fix graceful shutdown
	go func() {
		if err := e.Start(":" + strconv.Itoa(config.New().Service.DNA.Http.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Panic("shutting down the server, %v", err)
		}
	}()

	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Panic(err)
	}
}
