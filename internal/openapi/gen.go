package openapi

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -generate types -package openapi -o meetup-types.gen.go ../../openapi/meetup.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -generate server,spec -package openapi -o meetup-server.gen.go ../../openapi/meetup.yaml
//go:generate gofmt -s -w onboarding-server.gen.go meetup-types.gen.go
