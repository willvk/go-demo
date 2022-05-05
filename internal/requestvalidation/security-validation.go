package requestvalidation

import (
	"context"

	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
)

// NewSkipSecurityValidationOption returns request validation option
// that omits security validation. Security validation will be performed by
// API gateway, not by the middleware.
// Because API gateway is configured to use SigV4, it's not compatible
// with OpenAPI security scheme. Hence, omitting security validation
func NewSkipSecurityValidationOption() *oapimiddleware.Options {
	return &oapimiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
				// skip validation by returning nil
				return nil
			},
		},
	}
}
