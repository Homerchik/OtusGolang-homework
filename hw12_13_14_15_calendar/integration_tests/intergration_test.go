package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	suite.Suite
	ctx      context.Context
	cancel   context.CancelFunc
	HTTPAddr string
	MQUrl    string
	MQChan   <-chan models.EventMsg
}

const TimeFormat = "2006-01-02T15:04:05Z"

func (is *IntegrationSuite) PostReq(api string, data interface{}) *http.Response {
	jsonData, err := json.Marshal(data)
	is.Require().NoError(err)
	resp, err := http.Post(is.HTTPAddr+api, "application/json", bytes.NewBuffer(jsonData))
	is.Require().NoError(err)
	return resp
}

func (is *IntegrationSuite) GetReq(api string, params url.Values) *http.Response {
	var url string
	if len(params) > 0 {
		url = fmt.Sprintf("%s%s?%s", is.HTTPAddr, api, params.Encode())
	} else {
		url = fmt.Sprintf("%s%s", is.HTTPAddr, api)
	}
	resp, err := http.Get(url)
	is.Require().NoError(err)
	return resp
}

func (is *IntegrationSuite) DeleteReq(api string) *http.Response {
	url := fmt.Sprintf("%s%s", is.HTTPAddr, api)
	r, err := http.NewRequest("DELETE", url, nil)
	is.Require().NoError(err)
	client := &http.Client{}
	resp, err := client.Do(r)
	is.Require().NoError(err)
	return resp
}

func (is *IntegrationSuite) SetupSuite() {
	is.ctx, is.cancel = context.WithCancel(context.Background())
	is.HTTPAddr = "http://calendar:8080"
	is.MQUrl = "amqp://mq:5672"
	ch := make(chan models.EventMsg)
	is.MQChan = ch
	go NewReceiver(logger.New("DEBUG", "")).Run(is.ctx, is.MQUrl, "events", ch)
}

func (is *IntegrationSuite) TeardownSuite() {
	is.cancel()
}

func (is *IntegrationSuite) TestEventCRUDLifecycle() {
	startTime := time.Now().Unix() + 1000
	events := is.GenerateEventsForUser(uuid.New(), startTime, startTime+1, 1, 10)
	// CREATE
	event := is.AddEvents(events)[0]

	// GET CREATED
	resp := is.GetReq("/api/events/"+event.ID.String(), url.Values{})
	var respEvent models.Event
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	is.Require().NoError(err)
	is.Require().NoError(json.Unmarshal(body, &respEvent))
	is.Require().Equal(event, respEvent)

	// DELETE
	resp = is.DeleteReq("/api/events/" + event.ID.String())
	resp.Body.Close()
	is.Require().Equal(resp.StatusCode, http.StatusNoContent)

	// GET DELETED
	resp = is.GetReq("/api/events/"+event.ID.String(), url.Values{})
	resp.Body.Close()
	is.Require().Equal(resp.StatusCode, http.StatusNotFound)
}

func (is *IntegrationSuite) TestFullEventLifecycle() {
	afterNow := int64(100)
	startTime := time.Now().Unix() + afterNow
	event := is.GenerateEventsForUser(uuid.New(), startTime, startTime+1, 1, afterNow-2)[0]
	// CREATE
	resp := is.PostReq("/api/event", event)
	is.Require().Equal(resp.StatusCode, http.StatusCreated)
	id, err := uuid.Parse(strings.Split(resp.Header.Get("Location"), "/")[3])
	resp.Body.Close()
	is.Require().NoError(err)
	event.ID = id
	// GET
	resp = is.GetReq("/api/events/"+event.ID.String(), url.Values{})
	var respEvent models.Event
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	is.Require().NoError(err)
	is.Require().NoError(json.Unmarshal(body, &respEvent))
	is.Require().Equal(event, respEvent)
	// WAIT FOR MESSAGE TO FIRE
	select {
	case mqMessage := <-is.MQChan:
		is.Require().Equal(mqMessage.EventID, respEvent.ID)
	case <-time.After(10 * time.Second):
		is.FailNow("Message hasn't been received")
	}
}

func (is *IntegrationSuite) TestEventNotAddedNowGTStartDate() {
	startTime := time.Now().Unix() - 1000
	events := is.GenerateEventsForUser(uuid.New(), startTime, startTime+1, 1, 10)
	// CREATE
	resp := is.PostReq("/api/event", events[0])
	resp.Body.Close()
	is.Require().Equal(resp.StatusCode, http.StatusBadRequest)
}

func (is *IntegrationSuite) TestGetEventsForDay() {
	day := int64(60 * 60 * 24)
	startTime := time.Now().Unix() + day
	events := is.GenerateEventsForUser(uuid.New(), startTime, startTime+2*day, 3600, 10)
	// CREATE
	_ = is.AddEvents(events)

	// GET FOR A DAY
	params := url.Values{}
	params.Add("from", time.Unix(startTime, 0).Format(TimeFormat))
	params.Add("to", time.Unix(startTime+day, 0).Format(TimeFormat))
	resp := is.GetReq("/api/events", params)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var schedule map[string]models.Schedule
	is.Require().NoError(err)
	is.Require().NoError(json.Unmarshal(body, &schedule))
	for _, event := range schedule["data"] {
		is.Require().LessOrEqual(event.StartDate, startTime+day+3600)
	}
}

func (is *IntegrationSuite) AddEvents(events models.Schedule) models.Schedule {
	for i, event := range events {
		resp := is.PostReq("/api/event", event)
		resp.Body.Close()
		is.Require().Equal(resp.StatusCode, http.StatusCreated)
		id, err := uuid.Parse(strings.Split(resp.Header.Get("Location"), "/")[3])
		is.Require().NoError(err)
		events[i].ID = id
	}
	return events
}

func (is *IntegrationSuite) GenerateEventsForUser(
	userID uuid.UUID, startTS, endTS, interval, notifyBefore int64,
) models.Schedule {
	var schedule models.Schedule
	for i := startTS; i < endTS; i += interval {
		event := models.NewEvent(userID, "Some example title", "Just example", i, i+10, int(notifyBefore))
		schedule = append(schedule, event)
	}
	return schedule
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
