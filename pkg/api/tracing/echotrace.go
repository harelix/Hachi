package tracing

import (
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Middleware returns echo middleware which will trace incoming requests.
func Middleware(opts ...Option) echo.MiddlewareFunc {
	cfg := NewConfig(opts...)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			resource := request.Method + " " + c.Path()
			opts := []ddtrace.StartSpanOption{
				tracer.ServiceName(cfg.serviceName),
				tracer.ResourceName(resource),
				tracer.SpanType(ext.SpanTypeWeb),
				tracer.Tag(ext.HTTPMethod, request.Method),
				tracer.Tag(ext.HTTPURL, request.URL.Path),
				tracer.Measured(),
			}

			if !math.IsNaN(cfg.analyticsRate) {
				opts = append(opts, tracer.Tag(ext.EventSampleRate, cfg.analyticsRate))
			}
			if spanctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(request.Header)); err == nil {
				opts = append(opts, tracer.ChildOf(spanctx))
			}
			var finishOpts []tracer.FinishOption
			if cfg.noDebugStack {
				finishOpts = append(finishOpts, tracer.NoDebugStack())
			}
			span, ctx := tracer.StartSpanFromContext(request.Context(), "http.request", opts...)
			defer func() { span.Finish(finishOpts...) }()

			// pass the span through the request context
			c.SetRequest(request.WithContext(ctx))
			err := next(c)
			if err != nil {
				finishOpts = append(finishOpts, tracer.WithError(err))
				// invokes the registered HTTP error handler
				c.Error(err)
			}

			span.SetTag(ext.HTTPCode, strconv.Itoa(c.Response().Status))
			return err
		}
	}
}
