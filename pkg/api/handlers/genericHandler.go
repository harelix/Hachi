package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/internal"
	"github.com/rills-ai/Hachi/pkg/interpolator"
	"github.com/rills-ai/Hachi/pkg/messages"
	"github.com/rills-ai/Hachi/pkg/messaging"
	"github.com/rills-ai/Hachi/pkg/webhooks"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

func GenericHandler(c echo.Context, route config.RouteConfig) error {

	//webhooks.Construct().Notify("Yay!")
	selectors := route.Selectors
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
		Message:   body,
		Headers:   headers,
		Selectors: selectors,
		Route:     &route,
	}

	/*
		capsule string interpolation
	*/
	capsule = interpolateCapsuleValues(c, capsule)

	/*
		Dispatch capsule to subscribers
	*/
	response, err := DispatchCapsule(c, c.Request().Context(), capsule)

	if err != nil {
		log.Error("failed to dispatch request: %w", err)
		return echo.NewHTTPError(http.StatusServiceUnavailable, fmt.Errorf("failed to dispatch request: %w", err).Error())
	}
	return c.JSON(http.StatusOK, response)
}

func interpolateCapsuleValues(c echo.Context, capsule messages.Capsule) messages.Capsule {

	b, err := json.Marshal(capsule)
	if err != nil {
		log.Error("capsule message failed to interpolate, err: %w", err)
	}
	strCapsule := string(b)
	strCapsule = interpolator.InterpolateCapsuleValues(c, capsule.Route.IndexedInterpolationValues, strCapsule)

	var interpolatedCapsule messages.Capsule
	err = json.Unmarshal([]byte(strCapsule), &interpolatedCapsule)
	if err != nil {
		log.Error("capsule message failed to Unmarshal, err: %w", err)
	}
	return interpolatedCapsule
}

func DispatchCapsule(c echo.Context, ctx context.Context, capsule messages.Capsule) (messages.DefaultResponseMessage, error) {

	//remove internal instruction and handle a message by its remote sub-type!!!!!
	//checks for an internal instruction
	if capsule.Route.Remote.Internal != nil {
		directive := capsule.Route.Remote.Internal.Type
		response := internal.Exec(capsule, directive)
		m := messages.DefaultResponseMessage{
			Error:     false,
			Data:      response,
			Selectors: capsule.Selectors,
		}
		return m, nil
	}

	if capsule.Route.Remote.Webhook != nil {
		event := capsule.Route.Remote.Webhook.Event
		response := webhooks.Exec(capsule, event)
		m := messages.DefaultResponseMessage{
			Error:     false,
			Data:      response,
			Selectors: capsule.Selectors,
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
		Data:      HachiContext.PublishSuccessful,
		Selectors: capsule.Selectors,
	}
	return m, nil
}
