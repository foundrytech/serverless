package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"

	"github.com/mailgun/mailgun-go/v4"
)

func init() {
	functions.CloudEvent("VerifyEmail", verifyEmail)
}

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// NewUser is a struct representing a user account matching the JSON payload.
type NewUser struct {
	ID                      string      `json:"id"`
	FirstName               string      `json:"first_name"`
	LastName                string      `json:"last_name"`
	Username                string      `json:"username"`
	Verified                bool        `json:"verified"`
	VerificationToken       interface{} `json:"verification_token"`
	VerificationTokenCreated interface{} `json:"verification_token_created"`
	AccountCreated          string      `json:"account_created"`
	AccountUpdated          string      `json:"account_updated"`
}

// verifyAccount consumes a CloudEvent message and extracts the Pub/Sub message.
func verifyEmail(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	data := string(msg.Message.Data) // Automatically decoded from base64.
	var user NewUser
	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
    log.Printf("Error to unmarshal msg data: %v", err)
    return err
	}

	email := user.Username
	firstName := user.FirstName

	log.Printf("Email is: %s", email)
	log.Printf("First Name is: %s", firstName)


	// Generate verification token

	token, tokenCreated := generateVerificationToken()
	_ = tokenCreated
	_ = token

	// Update user data in db
	// db.UpdateUserVerificationToken(email, token, tokenCreated)


	// Send email to user
	sendMail(email, firstName, token)

	return nil
}

func generateVerificationToken() (string, time.Time) {
	token := uuid.New().String()
	tokenCreated := time.Now()
	return token, tokenCreated
}

func sendMail(email string, firstName string, token string) {
	domain := os.Getenv("DOMAIN_NAME")
	privateAPIKey := os.Getenv("MAILGUN_PRIVATE_API_KEY")

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(domain, privateAPIKey)

	sender := os.Getenv("SENDER")
	subject := os.Getenv("SUBJECT")
	recipient := email

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, "", recipient)

	url := fmt.Sprintf("http://%s:8080/v1/user/verify?token=%s", domain, token)

	body :=fmt.Sprintf(`
		<html>
			<body>
				<h1>Hi %s,</h1>
				<p style="font-size:30px;">Welcome to csye6225 ICU.</p>
				<p style="font-size:30px;">Please click the link below to verify your email address to get started. </p>
				<a href="%s" style="background-color: black; color: white; padding: 5x; text-align: center; text-decoration: none; display: inline-block; font-size: 16px; margin: 4px 2px; cursor: pointer; border: none;">
					Verify Email
				</a>
			</body>
		</html>`, firstName, url)

	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Send the message with a 5 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		fmt.Printf(`{"message": "Failed to send email: %s", "severity": "error"}`, err)
	}

	fmt.Printf(`{"message": "Email sent with resp %s, ID %s", "severity": "info"}`, resp, id)
}