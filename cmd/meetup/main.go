package main

import (
	"github.com/alecthomas/kong"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/willvk/go-demo/internal/app"
	"github.com/willvk/go-demo/internal/flags"
)

func main() {
	// Parse commandline flags / env vars
	var cmdFlags flags.Flags
	kong.Parse(&cmdFlags)

	// create middleware fields map to pass extra context to middlewares
	mwFields := map[string]interface{}{
		"lambda_container_id": app.GenerateUUID(), // Generate uuid for lambda container to track cold starts
	}
	mwFields = flags.GetFlagsMap(&cmdFlags, app.GetSourceDetailsMap(), mwFields)

	// Configure AWS Session and service clients
	// later
	// sess := session.Must(session.NewSession())

	// Setup API instance and its dependencies
	// Credential persistence client used by the Onboarding API to persist credentials
	//ddbClient := dynamodb.New(sess)
	//p := persistence.NewCredentialPersistence(ddbClient, cmdFlags.CredentialStoreName)
	//onboardingAPI := api.NewOnboardingAPI(p, tk, cmdFlags.CustomDomainName)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if cmdFlags.RunLocal {
		// Start the server (local run only) - useful for testing endpoints locally
		e.Logger.Fatal(e.Start(":3000"))
		return
	}

	// Use echoadapter to handle HTTP request from API gateway
	//echolambda := echoadapter.New(e)
	//handler := middlewareChain.Then(echolambda)
	//lambda.StartHandler(handler)
}
