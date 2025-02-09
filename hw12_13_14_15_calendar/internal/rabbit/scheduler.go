package rabbit

import (
	"context"
	"time"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/config"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/logic"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
)

type Scheduler struct {
	AMQPCon
	config  config.Scheduler
	storage models.Storage
}

const (
	pusherJob  = "pusher job"
	cleanerJob = "cleaner job"
)

func NewScheduler(config config.Scheduler, storage models.Storage, logger Logger) *Scheduler {
	return &Scheduler{
		AMQPCon: AMQPCon{name: "scheduler", logger: logger},
		config:  config,
		storage: storage,
	}
}

func (s *Scheduler) Run(ctx context.Context, amqpURL, queueName string) error {
	cleanerTicker := time.NewTicker(time.Duration(s.config.DeleteEvery) * time.Second)
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
		case <-cleanerTicker.C:
			if err := s.cleaner(); err != nil {
				s.logger.Error("%s: executing cleaner job: %v", cleanerJob, err)
			}
		}
	}
}

func (s *Scheduler) scanDB(ctx context.Context, out chan models.Event) error {
	defer close(out)
	for {
		curTS := time.Now().UTC().Unix()
		to := curTS + s.config.MaxNotifyBefore
		events, err := s.storage.GetEvents(curTS, to)
		if err != nil {
			s.logger.Error("%s: can't fetch events", pusherJob)
			return err
		}
		s.logger.Debug("%s: events %v", pusherJob, events)
		for _, event := range events {
			notifyTS := event.StartDate - int64(event.NotifyBefore)
			if curTS >= notifyTS && (curTS-notifyTS) <= s.config.ScanEvery {
				select {
				case <-ctx.Done():
				case out <- event:
					s.logger.Debug("%s: event has been sent %v", pusherJob, event.ID)
				}
			}
		}

		select {
		case <-ctx.Done():
			s.logger.Info("%s: closing scan-db goroutine", pusherJob)
			return nil
		case <-time.After(time.Duration(s.config.ScanEvery) * time.Second):
			s.logger.Info("%s: new fetch of events started", pusherJob)
		}
	}
}

func (s *Scheduler) cleaner() error {
	s.logger.Info("%s: delete goroutine started", cleanerJob)
	lastRunTime := time.Now().Unix()
	events, err := s.storage.GetEvents(0, lastRunTime-s.config.DeleteOlderThan)
	s.logger.Debug("%s: fetched events %v", cleanerJob, events)
	if err != nil {
		s.logger.Error("%s: can't fetch events: %v", cleanerJob, err.Error())
		return err
	}
	for _, e := range events {
		if err := s.storage.DeleteEvent(e.ID); err != nil {
			s.logger.Error("%s: can't delete event: %v", cleanerJob, e.ID)
		}
	}
	s.logger.Info(
		"%s: %d events older than %v successfully delete from DB",
		cleanerJob, len(events), s.config.DeleteOlderThan,
	)
	return nil
}
