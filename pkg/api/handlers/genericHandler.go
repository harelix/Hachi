package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rills-ai/Hachi/pkg/api/webhooks"
	"github.com/rills-ai/Hachi/pkg/helper"
	"github.com/rills-ai/Hachi/pkg/internal"
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

	webhooks.Construct().Notify("Yay!")
	//subjects := InterpolateContent(c, route)
	subjects := route.Subject
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
	interpolatedCapsule, err := interpolateCapsule(c, route, capsule)
	if err != nil {
		return err
	}
	fmt.Println(interpolatedCapsule)
	//todo add err handling
	response, err := DispatchCapsule(c, c.Request().Context(), interpolatedCapsule)

	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, fmt.Errorf("failed to dispatch request: %w", err).Error())
	}
	return c.JSON(http.StatusOK, response)
}

func DispatchCapsule(c echo.Context, ctx context.Context, capsule messages.Capsule) (map[string]any, error) {

	//remove internal instruction and handle a message by its remote sub-type!!!!!
	//checks for an internal instruction
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

	//main message/action relay/execution
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

//todo: interpolate route params for every field/member on our capsule (remote, subjects, headers,body, etc.)
func InterpolateContent(c echo.Context, route config.RouteConfig, content string) (interpolatedContent string) {

	for name, pattern := range route.IndexedInterpolationValues {
		value := c.Param(name)
		interpolatedContent = strings.TrimSpace(strings.Replace(content, pattern, value, -1))
	}

	return interpolatedContent
}

func interpolateCapsule(c echo.Context, route config.RouteConfig, capsule messages.Capsule) (messages.Capsule, error) {
	jsonCapsule, err := json.Marshal(capsule)
	if err != nil {
		return messages.Capsule{}, err
	}
	stringCapsule := string(jsonCapsule)

	interpolatedCapsule := InterpolateContent(c, route, stringCapsule)

	err = json.Unmarshal([]byte(interpolatedCapsule), &capsule)
	if err != nil {
		return messages.Capsule{}, err
	}

	return capsule, nil
}
