package routes

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/api/handlers"
	"github.com/rills-ai/Hachi/pkg/api/webhooks"
	"github.com/rills-ai/Hachi/pkg/config"
	"golang.org/x/exp/slices"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func RegisterRoutes(e *echo.Echo) error {

	apiVersion := strconv.Itoa(config.New().Service.DNA.API.Version)
	err := BindInstrumentationHandlers(e.Group("/"))
	if err != nil {
		return fmt.Errorf("failed to bind instrumentation: %w", err)
	}

	APIBaseRouteGroup := e.Group("/api/v" + apiVersion)
	err = BindRoutesFromConfiguration(APIBaseRouteGroup)
	if err != nil {
		return fmt.Errorf("failed to bind routes: %w", err)
	}

	/*
		Bind Webhooks
	*/
	webhooks.Construct().BindWebhooks(APIBaseRouteGroup)
	/*
		these methods should always be last (after routes registration or move them to )
		server discovery methods
	*/
	e.GET("/routes", func(routes []*echo.Route) echo.HandlerFunc {
		//_routes := routes
		return func(c echo.Context) error {
			return ListAllRoutes(c, routes)
		}
	}(e.Routes()))
	return nil
}

func ListAllRoutes(c echo.Context, routes []*echo.Route) error {
	return c.JSONPretty(http.StatusOK, routes, " ")
}

// BindRoutesFromConfiguration implements register and binds all routes declared in Hachi's configuration files
func BindRoutesFromConfiguration(group *echo.Group) error {

	internalStreams := config.New().Service.DNA.InternalTracts.Streams
	streams := append(internalStreams, config.New().Service.DNA.Tracts.Streams...)

	err := checkForRouteDuplicates(streams)
	if err != nil {
		return err
	}

	for _, stream := range streams {
		path := stream.Local

		if !(strings.HasPrefix(path, "/")) {
			path = "/" + path
		}

		//funny: just nice
		lcasedRoute := strings.ToLower(path)
		verb := strings.ToUpper(stream.Verb)
		//capture loop variables (route) in go closure
		//method := []string{route.Verb}
		var supportedMethods = []string{http.MethodPost, http.MethodGet}
		if !slices.Contains[string](supportedMethods, verb) {
			return errors.New("unsupported verb for route " + stream.Name)
		}

		registeredRoute := group.Add(verb, lcasedRoute,
			func(route config.RouteConfig) echo.HandlerFunc {
				return func(c echo.Context) error {
					return handlers.GenericHandler(c, DecorateRoute(route))
				}
			}(stream),
		)

		registeredRoute.Name = stream.Name
		fmt.Println(" config route register; name: " + registeredRoute.Name + " : " + registeredRoute.Path)
	}
	return nil
}

func checkForRouteDuplicates(streams []config.RouteConfig) error {
	routesMap := make(map[string]bool)
	for _, stream := range streams {
		routeDef := stream.Local
		if _, ok := routesMap[routeDef]; ok {
			return fmt.Errorf(HachiContext.DuplicateRouteDefinitionMessage, routeDef)
		}
		routesMap[routeDef] = true
	}
	return nil
}

func DecorateRoute(route config.RouteConfig) config.RouteConfig {
	return IndexRoutesInterpolationKeys(route)
}

var routeRegex = regexp.MustCompile("{{\\.((route)::(.*?))}}")

func IndexRoutesInterpolationKeys(route config.RouteConfig) config.RouteConfig {
	route.IndexedInterpolationValues = make(map[string]string)

	matches := routeRegex.FindAllStringSubmatch(strings.Join(route.Subject, " "), -1)
	for _, match := range matches {
		route.IndexedInterpolationValues[match[3]] = match[0]
	}
	return route
}
