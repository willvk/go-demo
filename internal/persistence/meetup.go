package persistence

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/honeycombio/beeline-go"
	"github.com/pkg/errors"
	"time"
)

const (
	meetupKey = "MEETUPID"
)

// MeetupPersistence defines interface for how to store generated bearer tokens
type MeetupPersistence interface {
	StoreMeetup(ctx context.Context, meetup *Meetup) (*Meetup, error)
	DeleteMeetup(ctx context.Context, meetupID string) (*Meetup, error)
}

type Meetup struct {
	MeetupID string

	AttendeeList []string

	// username of organiser.
	Organiser string

	// Meetup Date and Time
	PlannedDateTime time.Time

	// Meetup Reminder Date and Time
	RemindDateTime *time.Time
}

type MeetupRecord struct {
	//MEETUPID#{uuid}
	//USERID#{uuid} * future scope within the ddb
	PK              string
	SK              string
	MeetupID        string
	AttendeeList    []string
	Organiser       string
	PlannedDateTime string
	RemindDateTime  string
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

func (m MeetupStore) StoreMeetup(ctx context.Context, meetup *Meetup) (*Meetup, error) {
	ctx, span := beeline.StartSpan(ctx, "persistence.StoreMeetup")
	defer span.Send()

	newRecord := &MeetupRecord{
		PK:              fmt.Sprintf("%v#%v", meetupKey, meetup.MeetupID),
		SK:              meetup.PlannedDateTime.Format(time.RFC3339),
		MeetupID:        meetup.MeetupID,
		AttendeeList:    meetup.AttendeeList,
		Organiser:       meetup.Organiser,
		PlannedDateTime: meetup.PlannedDateTime.Format(time.RFC3339),
		RemindDateTime:  meetup.RemindDateTime.Format(time.RFC3339),
	}

	// ** Store the new meetup into the ddb
	itemRecord, errTo := dynamodbattribute.MarshalMap(newRecord)
	if errTo != nil {
		return nil, errors.Wrapf(errTo, "error cannot marshall items to store database object")
	}

	_, dbErr := m.ddb.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:      itemRecord,
		TableName: aws.String(m.databaseTableName),
	})
	if dbErr != nil {
		return nil, errors.Wrapf(dbErr, "error performing store operation on table")
	}
	return meetup, nil
}

func (m MeetupStore) DeleteMeetup(ctx context.Context, meetupID string) (*Meetup, error) {
	// delete by id then return deleted row
	return &Meetup{
		MeetupID:        "",
		AttendeeList:    nil,
		Organiser:       "",
		PlannedDateTime: time.Time{},
		RemindDateTime:  nil,
	}, nil
}
