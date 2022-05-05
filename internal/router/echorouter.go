package router

import (
	"fmt"
	"github.com/willvk/go-demo/internal/openapi"
	"github.com/willvk/go-demo/internal/requestvalidation"
	"os"
	"regexp"

	"github.com/labstack/echo/v4"

	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	sentrymiddleware "github.com/getsentry/sentry-go/echo"
	"github.com/honeycombio/beeline-go/wrappers/hnyecho"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// New creates an instance of Echo router
func New() *echo.Echo {
	// Setup OpenAPI spec
	swagger, err := openapi.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading Open API spec\n: %s", err)
		os.Exit(1)
	}
	// Clear out the servers array in the swagger spec to skip validating server names match.
	swagger.Servers = nil

	// Initialise echo, setup middleware and handlers
	e := echo.New()
	e.Pre(echomiddleware.RewriteWithConfig(echomiddleware.RewriteConfig{
		RegexRules: map[*regexp.Regexp]string{
			regexp.MustCompile(`^/\d*/(.*)`): "/$1",
		},
	}))
	e.Use(echomiddleware.Recover()) // recover from panic
	e.Use(echomiddleware.Gzip())    // enable compression
	e.Use(oapimiddleware.OapiRequestValidatorWithOptions(
		swagger,
		requestvalidation.NewSkipSecurityValidationOption()),
	) // validate input against OpenAPI schema

	e.Use(hnyecho.New().Middleware())                                                           // HoneyComb middleware
	e.Use(sentrymiddleware.New(sentrymiddleware.Options{Repanic: true, WaitForDelivery: true})) // Sentry middleware

	return e
}
