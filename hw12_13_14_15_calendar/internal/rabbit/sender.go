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
		AMQPCon: AMQPCon{name: "sender", logger: logger},
	}
}

func (s *Sender) Run(ctx context.Context, amqpURL, queueName string) error {
	if err := s.Connect(amqpURL, queueName, false); err != nil {
		return err
	}

	defer func() {
		if s.channel != nil {
			s.channel.Close()
		}
		if s.conn != nil {
			s.conn.Close()
		}
	}()

	consCh, err := s.channel.Consume(
		s.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.logger.Error("%s: can't create consume channel: %v: exiting...", err.Error())
		return err
	}
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("%s: sender finished", s.name)
			return nil
		case message := <-consCh:
			if err := s.ProcessMessage(message); err != nil {
				s.logger.Error("%s: can't process message: %v", s.name, err)
			}
		}
	}
}

func (s *Sender) ProcessMessage(message amqp.Delivery) error {
	var notification models.Notification
	if err := json.Unmarshal(message.Body, &notification); err != nil {
		return err
	}
	s.logger.Info("%s: notification received and logged: %v", s.name, notification)
	return nil
}
