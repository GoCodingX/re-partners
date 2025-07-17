package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	pkgerrors "github.com/GoCodingX/repartners/pkg/errors"
	"github.com/GoCodingX/repartners/pkg/gen/openapi"
	"github.com/GoCodingX/repartners/pkg/logger"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"
)

const timeoutInSeconds = 10

// TimeoutMiddleware sets a request timeout middleware with a custom
// "request timed out" message and a configurable timeout duration.
var TimeoutMiddleware = middleware.TimeoutWithConfig(middleware.TimeoutConfig{
	ErrorMessage: "request timed out",
	Timeout:      timeoutInSeconds * time.Second,
})

// OApiValidatorMiddleware creates an Echo middleware that validates requests
// against the provided OpenAPI spec and returns structured errors on validation failure.
func OApiValidatorMiddleware(swagger *openapi3.T) echo.MiddlewareFunc {
	return oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapimiddleware.Options{
		ErrorHandler:      pkgerrors.OApiErrorHandler(),
		MultiErrorHandler: pkgerrors.MultiErrorHandler(),
		Options: openapi3filter.Options{
			MultiError: true,
		},
	})
}

// CustomHTTPErrorHandler formats and sends consistent JSON error responses
// in openapi.ErrorResponse format for all errors returned by Echo handlers.
func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var (
		statusCode   int
		responseBody *openapi.ErrorResponse
	)

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		errorResponse, ok := echoErr.Message.(*openapi.ErrorResponse)
		if ok {
			statusCode = errorResponse.Code
			responseBody = errorResponse
		} else {
			statusCode = echoErr.Code
			msg, _ := echoErr.Message.(string)
			responseBody = pkgerrors.NewErrorResponse(statusCode, msg, nil)
		}
	} else {
		statusCode = http.StatusInternalServerError
		responseBody = pkgerrors.NewErrorResponse(statusCode, "something unexpected happened", nil)
	}

	if err := c.JSON(statusCode, responseBody); err != nil {
		logger.Error("failed to send error response", slog.String("err", err.Error()))
	}
}
