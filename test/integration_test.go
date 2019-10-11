package test

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/jslang/mailosaur-go/mailosaur"
	"github.com/stretchr/testify/require"
)

var (
	ServerID string
	APIKey   string
	SMTPPass string
)

func init() {
	gofakeit.Seed(0)
	ServerID = os.Getenv("MAILOSAUR_SERVER_ID")
	APIKey = os.Getenv("MAILOSAUR_API_KEY")
}

func sendMail(from string, to string, subject string, body string) error {
	msg := []byte(fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body))
	return smtp.SendMail("mailosaur.io:25", nil, from, []string{to}, msg)
}

func newEmail() string {
	return fmt.Sprintf("%s.%s@mailosaur.io", strings.ToLower(gofakeit.FirstName()), ServerID)
}

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

func TestMailosaur(t *testing.T) {

	testEmails := []*Email{
		{newEmail(), newEmail(), gofakeit.HackerVerb(), gofakeit.HackerPhrase()},
	}
	for _, email := range testEmails {
		if err := sendMail(email.From, email.To, email.Subject, email.Body); err != nil {
			panic(err)
		}
		log.Println("sent email", email)
	}

	c := mailosaur.NewClient(APIKey, ServerID)

	messages, err := c.ListMessages()
	require.NoError(t, err)
	require.Len(t, messages, len(testEmails))
	t.Log(messages[0].Id)
	message, err := c.GetMessage(messages[0].Id)
	require.NoError(t, err)
	require.NotEmpty(t, message)
	require.NoError(t, c.DeleteMessage(message.Id))
	require.NoError(t, c.DeleteMessages())
}
