package rabbit

import (
	"context"
	"time"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/logic"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
)

type Scheduler struct {
	AMQPCon
	storage         models.Storage
	scanEvery       int64
	maxNotifyBefore int64
	deleteEvery     int64
	deleteOlderThan int64
}

func NewScheduler(
	scanEvery, maxNotifyBefore, deleteEvery, deleteOlderThan int64,
	storage models.Storage, logger Logger,
) *Scheduler {
	return &Scheduler{
		AMQPCon:         AMQPCon{name: "scheduler", logger: logger},
		scanEvery:       scanEvery,
		maxNotifyBefore: maxNotifyBefore,
		deleteEvery:     deleteEvery,
		deleteOlderThan: deleteOlderThan,
		storage:         storage,
	}
}

func (s *Scheduler) Run(ctx context.Context, amqpURL, queueName string) error {
	if err := s.Connect(amqpURL, queueName, true); err != nil {
		return err
	}
	ch := make(chan models.Event)

	defer func() {
		if s.channel != nil {
			s.channel.Close()
		}
		if s.conn != nil {
			s.conn.Close()
		}
	}()

	go s.scanDB(ctx, ch)
	go s.cleaner(ctx)
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("%s: stopping scheduler", s.name)
			return nil
		case event := <-ch:
			notification := logic.BuildNotification(event)
			if err := s.PushJSON(notification); err != nil {
				s.logger.Error("%s: error sending notification: %v", s.name, err.Error())
				return err
			}
			s.logger.Info("%s: notification with id %v has been sent", s.name, notification.ID)
		}
	}
}

func (s *Scheduler) scanDB(ctx context.Context, out chan models.Event) error {
	name := "pusher job"
	defer close(out)
	for {
		curTS := time.Now().UTC().Unix()
		to := curTS + s.maxNotifyBefore
		events, err := s.storage.GetEvents(curTS, to)
		if err != nil {
			s.logger.Error("%s: can't fetch events", name)
			return err
		}
		s.logger.Debug("%s: events %v", name, events)
		for _, event := range events {
			notifyTS := event.StartDate - int64(event.NotifyBefore)
			if curTS >= notifyTS && (curTS-notifyTS) <= s.scanEvery {
				select {
				case <-ctx.Done():
				case out <- event:
					s.logger.Debug("%s: event has been sent %v", name, event.ID)
				}
			}
		}

		select {
		case <-ctx.Done():
			s.logger.Info("%s: closing scan-db goroutine", name)
			return nil
		case <-time.After(time.Duration(s.scanEvery) * time.Second):
			s.logger.Info("%s: new fetch of events started", name)
		}
	}
}

func (s *Scheduler) cleaner(ctx context.Context) {
	name := "cleaner job"
	s.logger.Info("%s: delete goroutine started", name)
	for {
		lastRunTime := time.Now().Unix()
		events, err := s.storage.GetEvents(0, lastRunTime-s.deleteOlderThan)
		s.logger.Debug("%s: fetched events %v", name, events)
		if err != nil {
			s.logger.Error("%s: can't fetch events: %v", name, err.Error())
		}
		for _, e := range events {
			if err := s.storage.DeleteEvent(e.ID); err != nil {
				s.logger.Error("%s: can't delete event: %v", name, e.ID)
			}
		}
		s.logger.Info("%s: %d events older than %v successfully delete from DB", name, len(events), s.deleteOlderThan)
		select {
		case <-ctx.Done():
			s.logger.Info("%s: delete job has closed", name)
		case <-time.After(time.Duration(s.deleteEvery) * time.Second):
			continue
		}
	}
}
