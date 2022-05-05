package persistence

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// MeetupPersistence defines interface for how to store generated bearer tokens
type MeetupPersistence interface {
}

// MeetupStore struct
type MeetupStore struct {
	ddb               dynamodbiface.DynamoDBAPI
	databaseTableName string
}

// NewMeetupPersistence creates a meetup persistence service
func NewMeetupPersistence(ddb dynamodbiface.DynamoDBAPI, databaseTableName string) MeetupPersistence {
	service := &MeetupStore{
		ddb:               ddb,
		databaseTableName: databaseTableName,
	}
	return service
}
