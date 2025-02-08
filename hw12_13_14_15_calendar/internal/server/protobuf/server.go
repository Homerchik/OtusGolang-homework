package protobuf

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/server/protobuf/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	storage    models.Storage
	logger     Logger
	grpcServer *grpc.Server
	pb.UnimplementedCalendarServer
}

type Logger interface {
	Info(format string, a ...any)
	Debug(format string, a ...any)
	Error(format string, a ...any)
}

var TimeFormat = "2006-01-02T15:04:05Z"

func NewServer(storage models.Storage, logger Logger) Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingServerInterceptor),
	)
	return Server{storage: storage, logger: logger, grpcServer: server}
}

func (s *Server) Start(ctx context.Context, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		s.logger.Error("can't start protobuf server: %w", err)
		return err
	}
	pb.RegisterCalendarServer(s.grpcServer, s)
	s.logger.Info("Protobuf server is running")
	if err := s.grpcServer.Serve(lis); err != nil {
		s.logger.Error("catched during serving: %w", err)
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) AddEvent(_ context.Context, req *pb.AddEventReq) (*pb.AddEventResp, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no data in request")
	}
	pbEvent := req.GetEvent()
	if pbEvent == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no event in request")
	}
	dbEvent := BuildEventFromPB(pbEvent)
	if dbEvent.UserID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}
	if err := s.storage.AddEvent(*dbEvent); err != nil {
		s.logger.Error("error adding event: %v", err)
		return nil, status.Errorf(codes.Internal, "error adding event: %v", err)
	}

	return &pb.AddEventResp{Id: dbEvent.ID.String()}, nil
}

func (s *Server) UpdateEvent(_ context.Context, req *pb.UpdateEventReq) (*pb.UpdateEventResp, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no data in request")
	}
	pbEvent := req.GetEvent()
	if pbEvent == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no event in request")
	}
	dbEvent := BuildEventFromPB(pbEvent)
	if err := s.storage.UpdateEvent(*dbEvent); err != nil {
		s.logger.Error("error updating event: %v", err)
		return nil, status.Errorf(codes.Internal, "error updating event: %v", err)
	}
	return &pb.UpdateEventResp{Id: dbEvent.ID.String()}, nil
}

func (s *Server) DeleteEvent(_ context.Context, req *pb.DeleteEventReq) (*pb.DeleteEventResp, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no data in request")
	}
	eventID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid event id")
	}
	if err := s.storage.DeleteEvent(eventID); err != nil {
		s.logger.Error("error deleting event: %v", err)
		return nil, status.Errorf(codes.Internal, "error deleting event: %v", err)
	}
	return &pb.DeleteEventResp{}, nil
}

func (s *Server) GetEvent(_ context.Context, req *pb.GetEventReq) (*pb.GetEventResp, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no data in request")
	}
	eventID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid event id")
	}
	_, event, err := s.storage.GetEventByID(eventID)
	if err != nil {
		s.logger.Error("error getting event: %v", err)
		return nil, status.Errorf(codes.Internal, "error getting event: %v", err)
	}
	return &pb.GetEventResp{Event: BuildPBFromEvent(event)}, nil
}

func (s *Server) GetEvents(_ context.Context, req *pb.GetEventsReq) (*pb.GetEventsResp, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no data in request")
	}
	fromDate, err := time.Parse(TimeFormat, req.GetFrom())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid from date")
	}
	toDate, err := time.Parse(TimeFormat, req.GetTo())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid to date")
	}
	events, err := s.storage.GetEvents(fromDate.Unix(), toDate.Unix())
	if err != nil {
		s.logger.Error("error getting events: %v", err)
		return nil, status.Errorf(codes.Internal, "error getting events: %v", err)
	}
	pbEvents := make([]*pb.Event, 0, len(events))
	for _, e := range events {
		pbEvents = append(pbEvents, BuildPBFromEvent(e))
	}
	return &pb.GetEventsResp{Events: pbEvents}, nil
}

func BuildEventFromPB(pbEvent *pb.Event) *models.Event {
	var (
		userID, eventID uuid.UUID
		err             error
	)
	userID, err = uuid.Parse(pbEvent.GetUserId())
	if err != nil {
		userID = uuid.Nil
	}
	eventID, err = uuid.Parse(pbEvent.GetId())
	if err != nil {
		eventID = uuid.Nil
	}
	return &models.Event{
		ID:           eventID,
		UserID:       userID,
		Title:        pbEvent.GetTitle(),
		Description:  pbEvent.GetDescription(),
		StartDate:    pbEvent.GetStartDate().Seconds,
		EndDate:      pbEvent.GetEndDate().Seconds,
		NotifyBefore: int(pbEvent.GetNotifyBefore()),
	}
}

func BuildPBFromEvent(event models.Event) *pb.Event {
	return &pb.Event{
		Id:           event.ID.String(),
		UserId:       event.UserID.String(),
		Title:        event.Title,
		Description:  event.Description,
		StartDate:    &timestamppb.Timestamp{Seconds: event.StartDate},
		EndDate:      &timestamppb.Timestamp{Seconds: event.EndDate},
		NotifyBefore: int64(event.NotifyBefore),
	}
}
