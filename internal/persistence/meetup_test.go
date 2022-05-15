package persistence

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/willvk/go-demo/mocks/awsmocks"
	"reflect"
	"testing"
	"time"
)

const (
	mockDDBTableName = "meetups-store"
)

var (
	mockPlannedDateTime = time.Now().Add(2 * time.Hour)
	mockOrganiserUN     = "willyvank"
	mockAttendeeList    = []string{"haley", "frank", "john"}
	mockMeetupID        = "abc123"
)

func Test_StoreMeetup(t *testing.T) {
	// assert := require.New(t)
	type fields struct {
		ddb               func(mock *awsmocks.MockDynamoDBAPI) *awsmocks.MockDynamoDBAPI
		databaseTableName string
	}
	type args struct {
		ctx   context.Context
		input *Meetup
		rec   *Meetup
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		errType string
	}{
		{
			name: "successfully stored meetup",
			fields: fields{
				ddb: func(mock *awsmocks.MockDynamoDBAPI) *awsmocks.MockDynamoDBAPI {
					newRecord := &MeetupRecord{
						PK:              fmt.Sprintf("%v#%v", meetupKey, mockMeetupID),
						SK:              mockPlannedDateTime.Format(time.RFC3339),
						MeetupID:        mockMeetupID,
						AttendeeList:    mockAttendeeList,
						Organiser:       mockOrganiserUN,
						PlannedDateTime: mockPlannedDateTime.Format(time.RFC3339),
						RemindDateTime:  "",
					}
					mock.EXPECT().PutItemWithContext(context.TODO(), newRecord).Return(nil, nil)
					return mock
				},
				databaseTableName: mockDDBTableName,
			},
			args: args{
				ctx: context.TODO(),
				input: &Meetup{
					MeetupID:        mockMeetupID,
					Organiser:       mockOrganiserUN,
					AttendeeList:    mockAttendeeList,
					PlannedDateTime: mockPlannedDateTime,
					RemindDateTime:  nil,
				},
				rec: &Meetup{
					MeetupID:        mockMeetupID,
					Organiser:       mockOrganiserUN,
					AttendeeList:    mockAttendeeList,
					PlannedDateTime: mockPlannedDateTime,
					RemindDateTime:  nil,
				},
			},
			wantErr: false,
			errType: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := awsmocks.NewMockDynamoDBAPI(ctrl)

			s := &MeetupStore{
				ddb:               tt.fields.ddb(mock),
				databaseTableName: tt.fields.databaseTableName,
			}
			got, err := s.StoreMeetup(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreCredential() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.args.rec) {
				t.Errorf("OAuthBearerToken \ngot = %v,\n want %v", got, tt.args.rec)
			}
		})
	}
}
