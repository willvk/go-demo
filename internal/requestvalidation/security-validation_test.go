package requestvalidation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSkipSecurityValidationOption(t *testing.T) {
	options := NewSkipSecurityValidationOption()

	assert.NotNil(t, options.Options.AuthenticationFunc)
	assert.Nil(t, options.Options.AuthenticationFunc(nil, nil))
}
