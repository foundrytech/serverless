package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"

	"github.com/jicodes/serverless/database"
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

	// Update user data in db
	// db.UpdateUserVerificationToken(email, token, tokenCreated)
	

	// Send email to user
	// email.SendVerificationEmail(email, firstName, token, tokenCreated)

	return nil
}

func generateVerificationToken() (string, time.Time) {
	token := uuid.New().String()
	tokenCreated := time.Now()
	return token, tokenCreated
}