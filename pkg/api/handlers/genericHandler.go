package handlers

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/internal"
	"github.com/rills-ai/Hachi/pkg/interpolator"
	"github.com/rills-ai/Hachi/pkg/messages"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

func GenericHandler(c echo.Context, route config.RouteConfig) error {

	//cfg := config.New()
	//todo: future task
	//cfg.Service.DNA.Stream.CircuitBreaker.Enabled
	//cfg.Service.DNA.Stream.Deduping.Enabled

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

	//fmt.Println(capsule.JSONFRomCapsule())
	/*
		capsule string interpolation
	*/
	capsule = interpolateCapsuleValues(c, capsule)

	/*
			======================================================
			Dispatch capsule to subscribers
		    ======================================================
	*/
	response, err := DispatchCapsule(c, c.Request().Context(), capsule)

	if err != nil {
		log.Error("failed to dispatch request: %w", err)
		return echo.NewHTTPError(http.StatusServiceUnavailable, fmt.Errorf("failed to dispatch request: %w", err).Error())
	}
	return c.JSON(http.StatusOK, response)
}

//tod: error handling
func interpolateCapsuleValues(c echo.Context, capsule messages.Capsule) messages.Capsule {

	capsuleString, err := capsule.JSONFRomCapsule()

	if err != nil {
		log.Error("capsule message failed to interpolate, err: %w", err)
		//tod: error handling
		return capsule
	}

	//convert route/path params to map
	var pathParamsDictionary = convertPathParamsToMap(c)

	//generic capsule interpolation method
	capsuleString, err = interpolator.InterpolateCapsuleValues(pathParamsDictionary,
		capsule.Route.IndexedInterpolationValues, capsuleString, true)

	if err != nil {
		log.Warning("capsule interpolation failed with err: %w", err)
	}

	interpolatedCapsule, err := messages.CapsuleFromJSON(capsuleString)

	if err != nil {
		log.Warning("capsule message failed to Unmarshal, err: %w", err)
		//todo: handle
		return capsule
	}
	return interpolatedCapsule
}

func convertPathParamsToMap(c echo.Context) map[string]string {

	var pathParams = make(map[string]string)
	var pNames = c.ParamNames()

	for i := 0; i < len(pNames); i++ {
		pathParams[pNames[i]] = c.ParamValues()[i]
	}
	return pathParams
}

func DispatchCapsule(c echo.Context, ctx context.Context, capsule messages.Capsule) (messages.DefaultResponseMessage, error) {
	//todo: implement Sync behaviour
	response, err := internal.ProcessCapsule(ctx, capsule)
	if err != nil {
		return messages.DefaultResponseMessage{
			Error: true,
		}, err
	}
	fmt.Println(response)

	m := messages.DefaultResponseMessage{
		Data:      HachiContext.PublishSuccessful,
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
