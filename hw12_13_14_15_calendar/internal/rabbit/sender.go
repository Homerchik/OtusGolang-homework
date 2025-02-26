package rabbit

import (
	"context"
	"encoding/json"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	AMQPCon
}

func NewSender(logger Logger) *Sender {
	return &Sender{
		AMQPCon: AMQPCon{Name: "sender", Logger: logger},
	}
}

func (s *Sender) Run(ctx context.Context, amqpURL, rcvQueue, pushQueue string) error {
	if err := s.Connect(amqpURL, true, pushQueue); err != nil {
		return err
	}

	defer func() {
		if s.Channel != nil {
			s.Channel.Close()
		}
		if s.Conn != nil {
			s.Conn.Close()
		}
	}()

	consCh, err := s.Channel.Consume(
		rcvQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.Logger.Error("%s: can't create consume channel: %v: exiting...", err.Error())
		return err
	}
	for {
		select {
		case <-ctx.Done():
			s.Logger.Info("%s: sender finished", s.Name)
			return nil
		case message := <-consCh:
			if err := s.ProcessMessage(message, pushQueue); err != nil {
				s.Logger.Error("%s: can't process message: %v", s.Name, err)
			}
		}
	}
}

func (s *Sender) ProcessMessage(message amqp.Delivery, pushQueue string) error {
	var notification models.Notification
	if err := json.Unmarshal(message.Body, &notification); err != nil {
		return err
	}
	s.Logger.Info("%s: notification received and logged: %v", s.Name, notification)
	if err := s.PushJSON(models.EventMsg{EventID: notification.ID, Status: "sent"}, pushQueue); err != nil {
		return err
	}
	s.Logger.Info("%s: notifications has been sent to push queue", s.Name)
	return nil
}
