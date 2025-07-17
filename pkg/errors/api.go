package errors

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/GoCodingX/repartners/pkg/gen/openapi"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ConvertEchoToApiError extracts and converts the payload inside an echo.HTTPError
// into a generated.ErrorResponse for consistent structured error handling across the API.
// Returns an error if the payload is not of type generated.ErrorResponse.
func ConvertEchoToApiError(err *echo.HTTPError) (*openapi.ErrorResponse, error) {
	errorResponse, ok := err.Message.(*openapi.ErrorResponse)
	if ok {
		return errorResponse, nil
	}

	msg, ok := err.Message.(string)
	if !ok {
		return nil, fmt.Errorf(
			"convertEchoToApiError: expected err.Message to be string, got %T (%v)",
			err.Message,
			err.Message,
		)
	}

	return NewErrorResponse(err.Code, msg, nil), nil
}

// OApiErrorHandler returns an Echo HTTP error handler that converts
// *echo.HTTPError instances produced by request validation failures into
// structured API error responses.
// Additionally, this ensures validation errors are returned consistently in generated.ErrorResponse
// format, allowing clients to receive detailed field-level validation feedback.
func OApiErrorHandler() func(c echo.Context, err *echo.HTTPError) error {
	return func(c echo.Context, err *echo.HTTPError) error {
		responsePayload, conversionErr := ConvertEchoToApiError(err)
		if conversionErr != nil {
			return fmt.Errorf("failed to conver to api error: %w", conversionErr)
		}

		return echo.NewHTTPError(err.Code, responsePayload)
	}
}

// openapi3.MultiError into a structured Echo HTTP 400 error.
func MultiErrorHandler() func(multiError openapi3.MultiError) *echo.HTTPError {
	return func(multiError openapi3.MultiError) *echo.HTTPError {
		return NewEchoBadRequestResponse(ParseSchemaErrors(multiError.Error()))
	}
}

// ParseSchemaErrors validation errors from a string into a slice of openapi.Detail.
func ParseSchemaErrors(input string) *[]openapi.Detail {
	var results []openapi.Detail

	// regex for query param errors
	reQuery := regexp.MustCompile(`parameter "([^"]+)" in query has an error: ([^\n|]+)`)

	// regex for body errors
	reBody := regexp.MustCompile(`Error at "/([^"]+)": ([^\n|]+)`)

	// process query param errors
	matchesQuery := reQuery.FindAllStringSubmatch(input, -1)
	for _, match := range matchesQuery {
		if len(match) == 3 {
			results = append(results, openapi.Detail{
				Field:   match[1],
				Message: strings.TrimSpace(match[2]),
			})
		}
	}

	// process body errors
	matchesBody := reBody.FindAllStringSubmatch(input, -1)
	for _, match := range matchesBody {
		if len(match) == 3 {
			rawField := match[1]
			field := strings.ReplaceAll(strings.TrimPrefix(rawField, "/"), "/", ".")
			results = append(results, openapi.Detail{
				Field:   field,
				Message: strings.TrimSpace(match[2]),
			})
		}
	}

	return &results
}

func NewEchoErrorResponse(statusCode int, message string, details *[]openapi.Detail) *echo.HTTPError {
	errorResponsePayload := NewErrorResponse(statusCode, message, details)

	return echo.NewHTTPError(errorResponsePayload.Code, errorResponsePayload)
}

func NewEchoBadRequestResponse(details *[]openapi.Detail) *echo.HTTPError {
	return NewEchoErrorResponse(http.StatusBadRequest, "request validation failed", details)
}

func NewErrorResponse(statusCode int, message string, details *[]openapi.Detail) *openapi.ErrorResponse {
	return &openapi.ErrorResponse{
		Code:    statusCode,
		Details: details,
		Message: message,
		Status:  http.StatusText(statusCode),
	}
}
