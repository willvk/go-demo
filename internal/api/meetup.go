package api

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"github.com/honeycombio/beeline-go"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/willvk/go-demo/internal/openapi"
	"github.com/willvk/go-demo/internal/persistence"
	"net/http"
	"time"
)

// NewMeetupAPI creates an instance of Onboarding API
func NewMeetupAPI(p persistence.MeetupPersistence, cdn string) *MeetupAPI {
	return &MeetupAPI{
		persistence: p,
		host:        cdn,
	}
}

type (
	// MeetupAPI implement meetup API operations
	MeetupAPI struct {
		persistence persistence.MeetupPersistence
		host        string
	}
)

func (m MeetupAPI) DeleteMeetup(ctx echo.Context) error {
	//TODO implement me
	return meetupError(ctx, http.StatusInternalServerError, "Error creating meetup. Reason: Not Yet Implemented")
}

func (m MeetupAPI) CreateMeetup(c echo.Context) error {
	ctx, span := beeline.StartSpan(c.Request().Context(), "meetups.CreateMeetup")
	defer span.Send()
	logger := log.With().Str("component", "meetup").Logger()
	loggctx := logger.WithContext(c.Request().Context())

	// process request body
	uuid, _ := uuid2.NewUUID()
	var createMeetupRequest openapi.CreateMeetupJSONBody
	if err := c.Bind(&createMeetupRequest); err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("MeetupID", uuid).Msg("error parsing request")
		return meetupError(c, http.StatusBadRequest, fmt.Sprintf("Error creating meetup. Reason: %s", err.Error()))
	}

	log.Ctx(loggctx).Info().Interface("MeetupID", uuid).Interface("Organiser", createMeetupRequest.Organiser).
		Msg(fmt.Sprintf("Creating Meetup for Date: %v", createMeetupRequest.PlannedDateTime.Format(time.RFC3339)))

	// if remind date is nil, it just gets passed through as nil
	_, persistenceError := m.persistence.StoreMeetup(ctx, &persistence.Meetup{
		MeetupID:        uuid.String(),
		AttendeeList:    createMeetupRequest.AttendeeList,
		Organiser:       createMeetupRequest.Organiser,
		PlannedDateTime: createMeetupRequest.PlannedDateTime,
		RemindDateTime:  createMeetupRequest.RemindDateTime,
	})
	if persistenceError != nil {
		log.Ctx(ctx).Error().Err(persistenceError).Interface("MeetupID", uuid).Msg("error storing meetup")
		return meetupError(c, http.StatusInternalServerError, "Error creating meetup. Internal Error")
	}

	return c.JSON(http.StatusCreated, &openapi.MeetupResponse{
		Meetup: openapi.Meetup{
			AttendeeList:    createMeetupRequest.AttendeeList,
			Organiser:       createMeetupRequest.Organiser,
			PlannedDateTime: createMeetupRequest.PlannedDateTime,
			RemindDateTime:  createMeetupRequest.RemindDateTime,
		},
		MeetupID: uuid.String(),
	})
}

func meetupError(c echo.Context, httpCode int, message string) error {
	return c.JSON(httpCode, &openapi.Error{
		Code:    int32(httpCode),
		Message: message,
	})
}
