package app

import (
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

// These are set at build time - see makefile
var (
	Name      = "unknown" // nolint:gochecknoglobals
	BuildDate = "unknown" // nolint:gochecknoglobals
	Commit    = "unknown" // nolint:gochecknoglobals
)

// SourceDetails is useful for debugging.
type SourceDetails struct {
	Name      string `mapstructure:"name"`
	BuildDate string `mapstructure:"-"`
	Commit    string `mapstructure:"commit"`
}

// New returns a SourceDetails
func New() SourceDetails {
	return SourceDetails{
		Name:      Name,
		BuildDate: BuildDate,
		Commit:    Commit,
	}
}

// GenerateUUID generates a UUID string. used for lambda cold start container tracking
func GenerateUUID() string {
	return uuid.New().String()
}

// GetSourceDetailsMap returns a map containing SourceDetails based on mapstructure keys
// It can be passed extra maps whose k:v's will also be added to the returned map
// SourceDetails k:v's trump passed in duplicates
func GetSourceDetailsMap(maps ...map[string]interface{}) map[string]interface{} {
	// resolve passed in maps to single map
	result := make(map[string]interface{})

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	// Decode Flags struct to json style
	jsonMap := make(map[string]interface{})

	err := mapstructure.Decode(New(), &jsonMap)
	if err != nil {
		log.Fatal().Err(err).Msg("App: failed to decode struct to map")
	}

	// Append decoded map to resolved map. Overwrites any duplicate keys
	for k, v := range jsonMap {
		result[k] = v
	}

	return result
}
