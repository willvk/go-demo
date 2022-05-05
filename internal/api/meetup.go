package api

import (
	"github.com/labstack/echo/v4"
	"github.com/willvk/go-demo/internal/openapi"
	"github.com/willvk/go-demo/internal/persistence"
	"net/http"
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

func (m MeetupAPI) CreateMeetup(ctx echo.Context) error {
	//TODO implement me
	return meetupError(ctx, http.StatusInternalServerError, "Error creating meetup. Reason: Not Yet Implemented")
}

func meetupError(c echo.Context, httpCode int, message string) error {
	return c.JSON(httpCode, &openapi.Error{
		Code:    int32(httpCode),
		Message: message,
	})
}
