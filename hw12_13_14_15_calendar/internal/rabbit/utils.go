package rabbit

import "fmt"

func BuildAMQPUrl(host string, port int, user, password string) string {
	if user == "" && password == "" {
		return fmt.Sprintf("amqp://%s:%d/", host, port)
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", user, password, host, port)
}
