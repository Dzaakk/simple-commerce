package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func Init(url string) (*Client, error) {
	for range 5 {
		conn, err := amqp.Dial(url)
		if err == nil {
			ch, err := conn.Channel()
			if err == nil {
				log.Print("Success connect to RabbitMQ")
				return &Client{
					conn:    conn,
					channel: ch,
				}, nil
			}
			conn.Close()
		}

		log.Print("RabbitMQ is not ready, retrying...")
		time.Sleep(5 * time.Second)
	}

	return nil, errors.New("failed to connect to RabbitMQ after multiple attempts")
}

func (c *Client) Channel() *amqp.Channel {
	if c == nil {
		return nil
	}
	return c.channel
}

func (c *Client) DeclareQueue(name string) (amqp.Queue, error) {
	if c == nil || c.channel == nil {
		return amqp.Queue{}, errors.New("rabbitmq client is not initialized")
	}

	return c.channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *Client) PublishJSON(ctx context.Context, queue string, payload any) error {
	if c == nil || c.channel == nil {
		return errors.New("rabbitmq client is not initialized")
	}

	if _, err := c.DeclareQueue(queue); err != nil {
		return err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return c.channel.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         body,
	})
}

func (c *Client) Close() error {
	if c == nil {
		return nil
	}

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
