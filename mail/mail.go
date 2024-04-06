package mail

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func Send(email string, firstName string, token string) {
	domain := os.Getenv("DOMAIN_NAME")
	privateAPIKey := os.Getenv("MAILGUN_PRIVATE_API_KEY")

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(domain, privateAPIKey)

	sender := os.Getenv("SENDER")
	subject := os.Getenv("SUBJECT")
	recipient := email

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, "", recipient)

	url := fmt.Sprintf("https://%s/v1/user/verify?token=%s", domain, token)

	body :=fmt.Sprintf(`
		<html>
			<body>
				<h1>Hi %s,</h1>
				<p style="font-size:20px;">Welcome to CSYE6225 ICU.</p>
				<p style="font-size:20px;">Please click the link below to verify your email address to get started. </p>
				<a href="%s" style="font-size:20px; background-color: #007bff; color: #fff; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">
					Verify Email Address
				</a>
				<p style="font-size:20px;">Thank you!</p>
				<p style="font-size:20px;">The CSYE6225 ICU Team</p>
			</body>
		</html>`, firstName, url)

	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Send the message with a 5 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Printf("Failed to send email: %v", err)
	}

	log.Printf("Email ID: %s sent, with resp: %s", id, resp)
}