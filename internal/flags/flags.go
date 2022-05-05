package flags

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Flags the flags/env vars passed to the service to configure it
// Limit these to items the user should pass in.
type Flags struct {
	RawEventLogging  bool   `help:"Enable raw event logging." env:"RAW_EVENT_LOGGING" mapstructure:"-"`
	LogLevel         string `help:"Set logging level" env:"LOG_LEVEL" mapstructure:"-"`
	Stage            string `help:"The development stage." env:"STAGE" mapstructure:"stage"`
	Branch           string `help:"The git branch this code originated." env:"BRANCH" mapstructure:"branch"`
	DependencyBranch string `help:"Git branch for downstream dependencies" env:"DEPENDENCY_BRANCH" mapstructure:"-"`
	AWSAccountID     string `help:"The aws account ID." env:"AWS_ACCOUNT_ID" mapstructure:"aws_account_id"`
	RunLocal         bool   `help:"Run HTTP server locally on port 3000" env:"RUN_LOCAL" mapstructure:"-"`
	SentryDSN        string `help:"The sentry DSN string for error logging." env:"SENTRY_DSN" mapstructure:"-"`
	MeetupStoreName  string `help:"AWS DynamoDB table name for Bearer Token Store" env:"MEETUP_TABLE_NAME" mapstructure:"-"`
	CustomDomainName string `help:"Custom Domain Name for the environment'" env:"DNS_NAME" mapstructure:"-"`
}

// ParseLogLevel parses value of LogLevel into a zerolog level
func ParseLogLevel(logLevel string) zerolog.Level {
	level := zerolog.InfoLevel

	if logLevel != "" {
		lowerLogLevel := strings.ToLower(logLevel)

		parsedLevel, err := zerolog.ParseLevel(lowerLogLevel)
		if err != nil {
			parsedLevel = zerolog.InfoLevel
			log.Warn().Msg("failed to parse Log level. Set level to info")
		}

		level = parsedLevel
	}

	return level
}

// GetFlagsMap returns a map containing Flags, based on mapstructure keys
// It can be passed extra maps whose k:v's will also be added to the returned map
// SourceDetails k:v's trump passed in duplicates
func GetFlagsMap(cmdFlags *Flags, maps ...map[string]interface{}) map[string]interface{} {
	// resolve passed in maps to single map
	result := make(map[string]interface{})

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	// Decode Flags struct to json style
	jsonMap := make(map[string]interface{})

	err := mapstructure.Decode(cmdFlags, &jsonMap)
	if err != nil {
		log.Fatal().Err(err).Msg("Flags: failed to decode struct to map")
	}

	// Append decoded map to resolved map. Overwrites any duplicate keys
	for k, v := range jsonMap {
		result[k] = v
	}

	return result
}
