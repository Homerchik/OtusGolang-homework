package rabbit

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPCon struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	name      string
	logger    Logger
}

type Logger interface {
	Info(format string, a ...any)
	Debug(format string, a ...any)
	Error(format string, a ...any)
}

func (a *AMQPCon) Connect(amqpURL string, queueName string, createQueue bool) error {
	a.queueName = queueName
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		a.logger.Error("%s: connecting to rabbitmq %v: %v", a.name, amqpURL, err.Error())
		return err
	}
	a.conn = conn
	a.logger.Info("%s: successfully connected to rabbitmq", a.name)

	a.channel, err = a.conn.Channel()
	if err != nil {
		a.logger.Error("%s: opening channel to rabbitmq: %v", a.name, err.Error())
		return err
	}
	a.logger.Info("%s: channel successfully opened", a.name)

	if createQueue {
		_, err = a.channel.QueueDeclare(
			a.queueName,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			a.logger.Error("%s: can't connect queue: %v", a.name, err)
			return err
		}
		a.logger.Info("%s: successfully created queue", a.name)
	}
	return nil
}

func (a *AMQPCon) PushJSON(v interface{}) error {
	name := "pusher job"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	a.logger.Debug("%s: message is ready to send: %v", name, v)
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	a.logger.Debug("%s: message has sent", name)
	return a.channel.PublishWithContext(
		ctx,
		"",
		a.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}
