package main

import (
	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/willvk/go-demo/internal/api"
	"github.com/willvk/go-demo/internal/app"
	"github.com/willvk/go-demo/internal/flags"
	"github.com/willvk/go-demo/internal/openapi"
	"github.com/willvk/go-demo/internal/persistence"
	"github.com/willvk/go-demo/internal/router"
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
	sess := session.Must(session.NewSession())

	// Setup API instance and its dependencies
	// Credential persistence client used by the Meetup API to persist meetups
	ddbClient := dynamodb.New(sess)
	p := persistence.NewMeetupPersistence(ddbClient, cmdFlags.MeetupStoreName)
	meetupAPI := api.NewMeetupAPI(p, cmdFlags.CustomDomainName)

	// Echo instance
	e := router.New()

	// Register route handlers
	openapi.RegisterHandlers(e, meetupAPI)

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
