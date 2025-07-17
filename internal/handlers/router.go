package handlers

import (
	"github.com/GoCodingX/repartners/internal/config"
	"github.com/GoCodingX/repartners/pkg/gen/openapi"
	middleware2 "github.com/GoCodingX/repartners/pkg/middleware"
	pkgmiddleware "github.com/GoCodingX/repartners/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(service *PacksService, swagger *openapi3.T, cfg *config.Config) (*echo.Echo, error) {
	e := echo.New()

	// Hide the default console output from echo
	e.HideBanner = true
	e.HidePort = true

	// middlewares
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// AllowOrigins: []string{"http://localhost:3000"},
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: []string{echo.OPTIONS, echo.DELETE},
		AllowHeaders: []string{echo.HeaderContentType},
	}))
	e.Use(middleware.RequestID())
	e.Use(middleware2.OApiValidatorMiddleware(swagger))
	e.Use(middleware2.TimeoutMiddleware)
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.HTTPErrorHandler = pkgmiddleware.CustomHTTPErrorHandler

	// register routes
	openapi.RegisterHandlers(e, service)

	return e, nil
}
