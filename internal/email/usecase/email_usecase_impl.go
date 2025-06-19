package usecase

import (
	"Dzaakk/simple-commerce/internal/email/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type EmailUsecCaseImpl struct {
}

func NewEmailUseCase() EmailUseCase {
	return &EmailUsecCaseImpl{}
}

func (e *EmailUsecCaseImpl) SendEmailActivation(ctx context.Context, data model.ActivationEmailReq) error {
	apiKey := os.Getenv("EMAIL_API_KEY")
	url := os.Getenv("EMAIL_API_URL")
	if apiKey == "" {
		return fmt.Errorf("API KEY is not set")
	}

	senderName := os.Getenv("SENDER_NAME")
	senderEmail := os.Getenv("SENDER_EMAIL")
	if senderName == "" || senderEmail == "" {
		return fmt.Errorf("sender information not set in environment variables")
	}

	emailReq := model.BaseEmailReq{
		Sender: model.Sender{
			Name:  senderName,
			Email: senderEmail,
		},
		To: []model.Recipient{{
			Email: data.Email,
			Name:  data.Username,
		},
		},
		Subject: "Activation Code",
		HTMLContent: fmt.Sprintf(`
			<html>
				<body>
					<h1>Welcome to Simple Commerce!</h1>
					<p>Hello %s,</p>
					<p>Thank you for signing up. Please use the following code to activate your account:</p>
					<h2 style="background-color: #f0f0f0; padding: 10px; text-align: center; font-size: 24px;">%s</h2>
					<p>This code will expire in 15 minutes.</p>
					<p>If you didn't request this code, please ignore this email.</p>
					<p>Best regards,<br>Simple Commerce Team</p>
				</body>
			</html>
		`, data.Username, data.ActivationCode),
	}

	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request to JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Response from Brevo API: %s", string(body))

	return nil
}
