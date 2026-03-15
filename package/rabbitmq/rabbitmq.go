package rabbitmq

import (
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Init(url string) (*amqp.Channel, error) {
	for range 5 {
		conn, err := amqp.Dial(url)
		if err == nil {
			ch, err := conn.Channel()
			if err == nil {
				log.Print("Success connect to RabbitMQ")
				return ch, nil
			}
			conn.Close()
		}

		log.Print("RabbitMQ is not ready, retrying...")
		time.Sleep(5 * time.Second)
	}

	return nil, errors.New("failed to connect to RabbitMQ after multiple attempts")
}
