package internalhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
	memorystorage "github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/suite"
)

type HTTPSuite struct {
	suite.Suite
	Storage models.Storage
	Server  Server
	mux     *http.ServeMux
}

var (
	hour   = int64(3600)
	future = time.Now().Add(time.Hour).UTC().Unix()
)

func (s *HTTPSuite) SetupSuite() {
	s.Storage = memorystorage.New()
	s.mux = NewServer("localhost", logger.New("debug", ""), s.Storage).mux
}

func (s *HTTPSuite) TestCreateHandler() {
	event := models.NewEvent(uuid.New(), "Event 1", "Best event", future+hour, future+2*hour, 60)
	bytesEvent, err := json.Marshal(event)
	s.Require().NoError(err)
	r, err := http.NewRequest(http.MethodPost, "http://localhost/api/event", bytes.NewBuffer(bytesEvent))
	s.Require().NoError(err)
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)
	s.Require().Equal(201, w.Code)
	id, err := uuid.Parse(strings.Split(w.Header().Get("Location"), "/")[3])
	s.Require().NoError(err)
	_, eventFromStorage, err := s.Storage.GetEventByID(id)
	s.Require().NoError(err)
	event.ID = id
	s.Equal(event, eventFromStorage)
}

func (s *HTTPSuite) TestUpdateHandler() {
	event := models.NewEvent(uuid.New(), "Event 1", "Best event", future+hour, future+2*hour, 60)
	s.Require().NoError(s.Storage.AddEvent(event))
	update, err := json.Marshal(map[string]interface{}{"title": "Event 2"})
	s.Require().NoError(err)
	event.Title = "Event 2"
	r, err := http.NewRequest(http.MethodPut, "http://localhost/api/events/"+event.ID.String(), bytes.NewBuffer(update))
	s.Require().NoError(err)
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)
	s.Require().Equal(204, w.Code)
	s.Require().Equal("/api/events/"+event.ID.String(), w.Header().Get("Location"))
	_, eventFromStorage, err := s.Storage.GetEventByID(event.ID)
	s.Require().NoError(err)
	s.Equal(event, eventFromStorage)
}

func (s *HTTPSuite) TestDeleteHandler() {
	event := models.NewEvent(uuid.New(), "Event 1", "Best event", future+hour, future+2*hour, 60)
	s.Require().NoError(s.Storage.AddEvent(event))
	r, err := http.NewRequest(http.MethodDelete, "http://localhost/api/events/"+event.ID.String(), nil)
	s.Require().NoError(err)
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)
	s.Require().Equal(204, w.Code)
	_, _, err = s.Storage.GetEventByID(event.ID)
	s.Require().Error(err)
}

func (s *HTTPSuite) TestGetEventHandler() {
	event := models.NewEvent(uuid.New(), "Event 1", "Best event", future+hour, future+2*hour, 60)
	s.Require().NoError(s.Storage.AddEvent(event))
	r, err := http.NewRequest(http.MethodGet, "http://localhost/api/events/"+event.ID.String(), nil)
	s.Require().NoError(err)
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)
	s.Require().Equal(200, w.Code)
	var eventFromResponse *models.Event
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &eventFromResponse))
	s.Equal(&event, eventFromResponse)
}

func (s *HTTPSuite) TestGetEventsHandler() {
	schedule := models.Schedule{
		models.NewEvent(uuid.New(), "Event 1", "Best event", future+8*hour, future+9*hour, 60),
		models.NewEvent(uuid.New(), "Event 1", "Best event", future+10*hour, future+11*hour, 60),
		models.NewEvent(uuid.New(), "Event 1", "Best event", future+12*hour, future+13*hour, 60),
		models.NewEvent(uuid.New(), "Event 1", "Best event", future+12*hour, future+13*hour, 60),
		models.NewEvent(uuid.New(), "Event 1", "Best event", future+12*hour, future+13*hour, 60),
	}
	for _, e := range schedule {
		s.Require().NoError(s.Storage.AddEvent(e))
	}
	fromTS := time.Now().UTC().Add(11 * time.Hour).Format(TimeFormat)
	toTS := time.Now().UTC().Add(14 * time.Hour).Format(TimeFormat)
	r, err := http.NewRequest(
		http.MethodGet, fmt.Sprintf("http://localhost/api/events?from=%s&to=%s", fromTS, toTS), nil,
	)
	s.Require().NoError(err)
	w := httptest.NewRecorder()
	s.mux.ServeHTTP(w, r)
	s.Require().Equal(200, w.Code)
	var eventFromResponse map[string]models.Schedule
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &eventFromResponse))
	s.Require().Equal(3, len(eventFromResponse["data"]))
}

func TestHTTPHandlers(t *testing.T) {
	suite.Run(t, new(HTTPSuite))
}
