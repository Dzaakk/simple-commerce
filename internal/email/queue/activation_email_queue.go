package queue

import (
	emailModel "Dzaakk/simple-commerce/internal/email/model"
	emailService "Dzaakk/simple-commerce/internal/email/service"
	"Dzaakk/simple-commerce/package/logging"
	"Dzaakk/simple-commerce/package/rabbitmq"
	"context"
	"encoding/json"
	"errors"
	"log"
)

const ActivationEmailQueue = "email.activation.verify"

type ActivationEmailJob struct {
	Email          string `json:"email"`
	Username       string `json:"username"`
	ActivationLink string `json:"activation_link"`
}

type RabbitPublisher struct {
	Client *rabbitmq.Client
	Logger *logging.Logger
}

func NewRabbitPublisher(client *rabbitmq.Client) *RabbitPublisher {
	return &RabbitPublisher{
		Client: client,
		Logger: logging.NewLogger("email", "activation_email_publisher"),
	}
}

func (p *RabbitPublisher) PublishVerificationEmail(ctx context.Context, req emailModel.VerificationEmailReq) error {
	if p == nil || p.Client == nil {
		return errors.New("rabbitmq publisher is not initialized")
	}

	job := ActivationEmailJob{
		Email:          req.Email,
		Username:       req.Username,
		ActivationLink: req.ActivationLink,
	}

	if err := p.Client.PublishJSON(ctx, ActivationEmailQueue, job); err != nil {
		p.Logger.Error(ctx, "activation_email_publish_failed", map[string]interface{}{
			"queue": ActivationEmailQueue,
		})
		return err
	}

	p.Logger.Info(ctx, "activation_email_queued", map[string]interface{}{
		"queue": ActivationEmailQueue,
	})
	return nil
}

func StartActivationEmailConsumer(client *rabbitmq.Client, svc emailService.EmailService) error {
	if client == nil {
		return errors.New("rabbitmq client is not initialized")
	}
	if svc == nil {
		return errors.New("email service is not initialized")
	}

	if _, err := client.DeclareQueue(ActivationEmailQueue); err != nil {
		return err
	}

	msgs, err := client.Channel().Consume(
		ActivationEmailQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	logger := logging.NewLogger("email", "activation_email_consumer")

	go func() {
		for msg := range msgs {
			var job ActivationEmailJob
			if err := json.Unmarshal(msg.Body, &job); err != nil {
				logger.Error(context.Background(), "activation_email_invalid_payload", map[string]interface{}{
					"queue": ActivationEmailQueue,
				})
				_ = msg.Reject(false)
				continue
			}

			err := svc.SendEmailVerification(context.Background(), emailModel.VerificationEmailReq{
				Email:          job.Email,
				Username:       job.Username,
				ActivationLink: job.ActivationLink,
			})
			if err != nil {
				logger.Error(context.Background(), "activation_email_send_failed", map[string]interface{}{
					"queue": ActivationEmailQueue,
				})

				if msg.Redelivered {
					_ = msg.Reject(false)
				} else {
					_ = msg.Nack(false, true)
				}
				continue
			}

			logger.Info(context.Background(), "activation_email_sent", map[string]interface{}{
				"queue": ActivationEmailQueue,
			})

			if err := msg.Ack(false); err != nil {
				log.Printf("failed to ack activation email message: %v", err)
			}
		}
	}()

	return nil
}
