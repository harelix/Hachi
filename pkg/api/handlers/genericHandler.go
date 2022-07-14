package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rills-ai/Hachi/pkg/internal"
	"github.com/rills-ai/Hachi/pkg/webhooks"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/messages"
	"github.com/rills-ai/Hachi/pkg/messaging"
)

func GenericHandler(c echo.Context, route config.RouteConfig) error {

	//webhooks.Construct().Notify("Yay!")
	subjects := InterpolateRoutingKeyFromRouteParams(c, route)
	headers := c.Request().Header
	body := route.Payload

	/*
		comment: no need for other HTTP verbs:
		-	Post and Get are symbolic for variant message and const message
	*/
	if strings.ToUpper(route.Verb) == http.MethodPost {
		b, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "failed to read request body")
		} else if len(b) == 0 {
			return echo.NewHTTPError(http.StatusExpectationFailed, HachiContext.ErrBodyEmpty.Error())
		}
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

	b, err := json.Marshal(capsule)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))

	//todo add err handling
	response, err := DispatchCapsule(c, c.Request().Context(), capsule)

	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, fmt.Errorf("failed to dispatch request: %w", err).Error())
	}
	return c.JSON(http.StatusOK, response)
}

func DispatchCapsule(c echo.Context, ctx context.Context, capsule messages.Capsule) (messages.DefaultResponseMessage, error) {

	//remove internal instruction and handle a message by its remote sub-type!!!!!
	//checks for an internal instruction
	if capsule.Route.Remote.Internal != nil {
		directive := capsule.Route.Remote.Internal.Type
		response := internal.Exec(capsule, directive)
		m := messages.DefaultResponseMessage{
			Error: false,
			Data:  response,
			Path:  capsule.Subject,
		}
		return m, nil
	}

	if capsule.Route.Remote.Webhook != nil {
		event := capsule.Route.Remote.Webhook.Event
		response := webhooks.Exec(capsule, event)
		m := messages.DefaultResponseMessage{
			Error: false,
			Data:  response,
			Path:  capsule.Subject,
		}
		return m, nil
	}
	/*
			directive := helper.CollectionFunc[string](capsule.Subject,
				func(value string) bool {
					return strings.Contains(value, "__internal__")
				})
		if directive != "" {
			//internal actions execution
			response := internal.Exec(capsule, directive)
			m := map[string]any{
				"data": response,
				"path": capsule.Subject,
			}
			return m, nil
		}
	*/
	//main message/action relay/execution
	err := messaging.Get().Publish(ctx, capsule)
	if err != nil {
		return messages.DefaultResponseMessage{
			Error: true,
		}, err
	}

	m := messages.DefaultResponseMessage{
		Data: HachiContext.PublishSuccessful,
		Path: capsule.Subject,
	}
	return m, nil
}

//todo: interpolate route params for every field/member on our capsule (remote, subjects, headers,body, etc.)
func InterpolateRoutingKeyFromRouteParams(c echo.Context, route config.RouteConfig) []string {

	for name, pattern := range route.IndexedInterpolationValues {
		value := c.Param(name)
		for idx, Avalue := range route.Selectors {
			route.Selectors[idx] = strings.TrimSpace(strings.Replace(Avalue, pattern, value, -1))
		}
	}
	return route.Selectors
}
