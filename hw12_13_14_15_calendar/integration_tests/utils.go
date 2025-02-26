package integrationtests

import (
	"context"
	"encoding/json"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/rabbit"
)

type Receiver struct {
	rabbit.AMQPCon
}

func NewReceiver(logger rabbit.Logger) *Receiver {
	return &Receiver{
		AMQPCon: rabbit.AMQPCon{Name: "receiver", Logger: logger},
	}
}

func (r *Receiver) Run(ctx context.Context, amqpURL, queueName string, out chan<- models.EventMsg) error {
	if err := r.Connect(amqpURL, false); err != nil {
		return err
	}

	defer func() {
		if r.Channel != nil {
			r.Channel.Close()
		}
		if r.Conn != nil {
			r.Conn.Close()
		}
		close(out)
	}()

	consCh, err := r.Channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case message := <-consCh:
			var queueMsg models.EventMsg
			if err := json.Unmarshal(message.Body, &queueMsg); err != nil {
				return err
			}
			out <- queueMsg
		}
	}
}
