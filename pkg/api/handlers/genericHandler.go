package handlers

import (
	"context"
	"fmt"
	"github.com/gilbarco-ai/Hachi/pkg/helper"
	"github.com/gilbarco-ai/Hachi/pkg/internal"
	"io"
	"net/http"
	"strings"

	HachiContext "github.com/gilbarco-ai/Hachi/pkg"
	"github.com/gilbarco-ai/Hachi/pkg/config"
	"github.com/gilbarco-ai/Hachi/pkg/messages"
	"github.com/gilbarco-ai/Hachi/pkg/messaging"
	"github.com/labstack/echo/v4"
)

func HachiGenericHandler(c echo.Context, route config.RouteConfig) error {

	subjects := InterpolateRoutingKeyFromRouteParams(c, route)
	headers := c.Request().Header
	body := route.Payload

	//todo: no need for other HTTP verbs / Post and Get are symbolic for variant message and const message
	if strings.ToUpper(route.Verb) == http.MethodPost {
		b, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "failed to read request body")
		} else if len(b) == 0 {
			return echo.NewHTTPError(http.StatusExpectationFailed, HachiContext.ErrBodyEmpty.Error())
		}
		// TODO what if payload is not empty
		body = string(b)
	}

	for k, v := range route.Headers {
		headers.Set(k, strings.Join(v, ";"))
	}

	capsule := messages.Capsule{
		Message: body,
		Headers: headers,
		Subject: subjects,
		Route:   &route,
	}

	//todo add err handling
	response, err := DispatchCapsule(c, c.Request().Context(), capsule)

	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, fmt.Errorf("failed to dispatch request: %w", err).Error())
	}
	return c.JSON(http.StatusOK, response)
}

func DispatchCapsule(c echo.Context, ctx context.Context, capsule messages.Capsule) (map[string]any, error) {

	//checks for an internal instruction
	directive := helper.CollectionFunc[string](capsule.Subject,
		func(value string) bool {
			return strings.Contains(value, "__internal__")
		})

	if directive != "" {
		response := internal.Exec(capsule, directive)
		m := map[string]any{
			"data": response,
			"path": capsule.Subject,
		}
		return m, nil
	} else {
		err := messaging.Get().Publish(ctx, capsule)
		if err != nil {
			return nil, err
		}
		m := map[string]any{
			"data": HachiContext.PublishSuccessful,
			"path": capsule.Subject,
		}
		return m, nil
	}
}

func InterpolateRoutingKeyFromRouteParams(c echo.Context, route config.RouteConfig) []string {

	for name, pattern := range route.IndexedInterpolationValues {
		value := c.Param(name)
		for idx, Avalue := range route.Subject {
			route.Subject[idx] = strings.TrimSpace(strings.Replace(Avalue, pattern, value, -1))
		}
	}
	return route.Subject
}
