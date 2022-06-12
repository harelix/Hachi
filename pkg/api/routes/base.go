package routes

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gilbarco-ai/Hachi/pkg/api/handlers"
	"github.com/gilbarco-ai/Hachi/pkg/config"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

func RegisterRoutes(e *echo.Echo) error {

	apiVersion := strconv.Itoa(config.New().Service.DNA.API.Version)

	err := BindInstrumentationHandlers(e.Group("/"))
	if err != nil {
		return fmt.Errorf("failed to bind instrumentation: %w", err)
	}

	versionedAPI := e.Group("/api/v" + apiVersion)
	err = BindRoutesFromConfiguration(versionedAPI)
	if err != nil {
		return fmt.Errorf("failed to bind routes: %w", err)
	}
	//these methods should always be last (after routes registration or move them to )
	//server discovery methods
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

	streams := config.New().Service.DNA.Tracts.Streams

	for _, stream := range streams {
		path := stream.Local

		if !(strings.HasPrefix(path, "/")) {
			path = "/" + path
		}

		//just nice
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
					return handlers.HachiGenericHandler(c, DecorateRoute(route))
				}
			}(stream),
		)

		registeredRoute.Name = stream.Name
		fmt.Println(" config route register; name: " + registeredRoute.Name + " : " + registeredRoute.Path)
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
