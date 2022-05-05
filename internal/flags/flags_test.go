package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rs/zerolog"
)

// Tests

func TestParseLogLevel(t *testing.T) {
	tests := map[string]struct {
		level string
		want  zerolog.Level
	}{
		"should return valid level if lowercase": {
			level: "debug",
			want:  zerolog.DebugLevel,
		},
		"should return valid level if uppercase": {
			level: "ERROR",
			want:  zerolog.ErrorLevel,
		},
		"should return info level if empty": {
			level: "",
			want:  zerolog.InfoLevel,
		},
		"should return info level if invalid": {
			level: "garbage",
			want:  zerolog.InfoLevel,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// run
			got := ParseLogLevel(tt.level)

			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFlagsMap(t *testing.T) {
	tests := []struct {
		name     string
		input    Flags
		expected map[string]interface{}
	}{
		{
			name: "flags successfully converted to map",
			input: Flags{
				AWSAccountID:     "123456789012",
				Stage:            "dev",
				Branch:           "master",
				DependencyBranch: "unittest",
			},
			expected: map[string]interface{}{
				"aws_account_id": "123456789012",
				"stage":          "dev",
				"branch":         "master",
			},
		},
		{
			name: "flags correctly excludes values in map",
			input: Flags{
				AWSAccountID:     "123456789012",
				Stage:            "dev",
				Branch:           "master",
				SentryDSN:        "sentry.io/some-random-dsn/123456789",
				RunLocal:         true,
				DependencyBranch: "unittest",
			},
			expected: map[string]interface{}{
				"aws_account_id": "123456789012",
				"stage":          "dev",
				"branch":         "master",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFlagsMap(&tt.input)
			assert.EqualValues(t, tt.expected, result)
		})
	}
}
