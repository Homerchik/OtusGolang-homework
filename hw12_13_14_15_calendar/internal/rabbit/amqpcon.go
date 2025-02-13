package rabbit

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPCon struct {
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	Name      string
	Logger    Logger
}

type Logger interface {
	Info(format string, a ...any)
	Debug(format string, a ...any)
	Error(format string, a ...any)
}

func (a *AMQPCon) Connect(amqpURL string, createQueue bool, queues ...string) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		a.Logger.Error("%s: connecting to rabbitmq %v: %v", a.Name, amqpURL, err.Error())
		return err
	}
	a.Conn = conn
	a.Logger.Info("%s: successfully connected to rabbitmq", a.Name)

	a.Channel, err = a.Conn.Channel()
	if err != nil {
		a.Logger.Error("%s: opening channel to rabbitmq: %v", a.Name, err.Error())
		return err
	}
	a.Logger.Info("%s: channel successfully opened", a.Name)

	if createQueue {
		for _, q := range queues {
			_, err = a.Channel.QueueDeclare(
				q,
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				a.Logger.Error("%s: can't create queue: %v", a.Name, err)
				return err
			}
			a.Logger.Info("%s: successfully created queue", a.Name)
		}
	}
	return nil
}

func (a *AMQPCon) PushJSON(v interface{}, queueName string) error {
	name := "pusher job"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	a.Logger.Debug("%s: message is ready to send: %v", name, v)
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	a.Logger.Debug("%s: message has sent", name)
	return a.Channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}
