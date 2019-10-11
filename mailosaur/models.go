package mailosaur

import (
	"encoding/json"
	"time"
)

// baseMessage provides fields common between the full and summary versions of Message objects
type baseMessage struct {
	Id       string              `json:"id"`
	Server   string              `json:"server"`
	From     []map[string]string `json:"from"`
	To       []map[string]string `json:"to"`
	CC       []map[string]string `json:"cc"`
	BCC      []map[string]string `json:"bcc"`
	Received time.Time           `json:"received"`
	Subject  string              `json:"subject"`
	Summary  string              `json:"summary"`
}

// Message objects represent an email or SMS received by Mailosaur and contain all the data you might need to perform
// any number of manual or automated tests.
type Message struct {
	baseMessage

	// TODO: implement explictly defined structs for below
	Attachments json.RawMessage `json:"attachments"`
	HTML        json.RawMessage `json:"html"`
	Text        json.RawMessage `json:"text"`
	Metadata    json.RawMessage `json:"metadata"`
	HATEOSLinks json.RawMessage `json:"hateosLinks"`
}

// MessageSummary objects represent a summarized email or SMS received by Mailosaur.
type MessageSummary struct {
	baseMessage

	Attachments int `json:"attachments"`
}
