package protobuf

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/server/protobuf/pb"
	memorystorage "github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCSuite struct {
	suite.Suite
	Storage models.Storage
	Server  Server
	Client  pb.CalendarClient
	cancel  context.CancelFunc
	ctx     context.Context
}

var (
	hour   = int64(3600)
	future = time.Now().Add(time.Hour).UTC().Unix()
)

func (s *GRPCSuite) SetupSuite() {
	s.Storage = memorystorage.New()
	s.ctx, s.cancel = context.WithCancel(context.Background())

	port := 18090
	s.Server = NewServer(s.Storage, logger.New("debug", ""))
	go s.Server.Start(s.ctx, port)
	conn, err := grpc.NewClient(fmt.Sprintf(":%v", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Client = pb.NewCalendarClient(conn)
}

func (s *GRPCSuite) TeardownSuite() {
	s.cancel()
}

// func (s *GRPCSuite) AfterTest(suiteName, testName string) {
// 	s.Storage = memorystorage.New()
// }

func (s *GRPCSuite) TestCreateHandler() {
	event := models.NewEvent(uuid.New(), "Event 1", "Best event", future+hour, future+2*hour, 60)
	resp, err := s.Client.AddEvent(s.ctx, &pb.AddEventReq{Event: BuildPBFromEvent(event)})
	s.Require().NoError(err)
	id, err := uuid.Parse(resp.Id)
	s.Require().NoError(err)
	_, eventFromStorage, err := s.Storage.GetEventByID(id)
	s.Require().NoError(err)
	event.ID = id
	s.Equal(event, eventFromStorage)
}

func (s *GRPCSuite) TestUpdateHandler() {
	event := models.NewEvent(uuid.New(), "Event 1", "Best event", future+hour, future+2*hour, 60)
	s.Require().NoError(s.Storage.AddEvent(event))
	event.Title = "Event 2"
	resp, err := s.Client.UpdateEvent(s.ctx, &pb.UpdateEventReq{Event: BuildPBFromEvent(event)})
	s.Require().NoError(err)
	id, err := uuid.Parse(resp.Id)
	s.Require().NoError(err)
	_, eventFromStorage, err := s.Storage.GetEventByID(id)
	s.Require().NoError(err)
	s.Equal(event, eventFromStorage)
}

func (s *GRPCSuite) TestDeleteHandler() {
	event := models.NewEvent(uuid.New(), "Event 1", "Best event", future+hour, future+2*hour, 60)
	s.Require().NoError(s.Storage.AddEvent(event))
	_, err := s.Client.DeleteEvent(s.ctx, &pb.DeleteEventReq{Id: event.ID.String()})
	s.Require().NoError(err)
	_, _, err = s.Storage.GetEventByID(event.ID)
	s.Require().Error(err)
}

func (s *GRPCSuite) TestGetEventHandler() {
	event := models.NewEvent(uuid.New(), "Event 1", "Best event", future+hour, future+2*hour, 60)
	s.Require().NoError(s.Storage.AddEvent(event))
	resp, err := s.Client.GetEvent(s.ctx, &pb.GetEventReq{Id: event.ID.String()})
	s.Require().NoError(err)
	pbEvent := BuildEventFromPB(resp.Event)
	s.Equal(&event, pbEvent)
}

func (s *GRPCSuite) TestGetEventsHandler() {
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
	FromTS := time.Now().UTC().Add(11 * time.Hour).Format(TimeFormat)
	toTS := time.Now().UTC().Add(14 * time.Hour).Format(TimeFormat)
	resp, err := s.Client.GetEvents(s.ctx, &pb.GetEventsReq{From: FromTS, To: toTS})
	s.Require().NoError(err)
	var eventFromPB models.Schedule
	for _, e := range resp.Events {
		eventFromPB = append(eventFromPB, *BuildEventFromPB(e))
	}
	s.Require().Equal(3, len(eventFromPB))
}

func TestGRPCServer(t *testing.T) {
	suite.Run(t, new(GRPCSuite))
}
