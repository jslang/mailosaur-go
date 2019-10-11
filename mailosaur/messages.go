package mailosaur

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// GetMessage retrieves the detail for a single email message.
func (c *Client) GetMessage(messageID string) (*Message, error) {
	httpResp, err := c.call(http.MethodGet, "messages/"+messageID, nil, nil)
	if err != nil {
		return nil, err
	}
	var msg Message
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	return &msg, json.Unmarshal(body, &msg)
}

// DeleteMessage permanently deletes a message.
func (c *Client) DeleteMessage(messageID string) error {
	_, err := c.call(http.MethodDelete, "messages/"+messageID, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

type (
	ListMessagesRequest struct {
		Page          int
		ItemsPerPage  int
		ReceivedAfter time.Time
	}
)

type listMessagesOption func(map[string]interface{})

// SetListMessagesPage sets the page of results to request
func SetListMessagesPage(page int) listMessagesOption {
	return func(data map[string]interface{}) {
		data["page"] = page
	}
}

// SetListMessagesItemsPerPage sets the number of items to display per result page
func SetListMessagesItemsPerPage(itemsPerPage int) listMessagesOption {
	return func(data map[string]interface{}) {
		data["itemsPerPage"] = itemsPerPage
	}
}

// SetListMessagesReceivedAfter sets the time to filter messages by
func SetListMessagesReceivedAfter(receivedAfter time.Time) listMessagesOption {
	return func(data map[string]interface{}) {
		data["receivedAfter"] = receivedAfter.Format(time.RFC3339)
	}
}

// ListMessages returns a list of your messages in summary form.
func (c *Client) ListMessages(options ...listMessagesOption) ([]*MessageSummary, error) {
	queryParams := map[string]interface{}{
		"server": c.serverID,
	}
	for _, opt := range options {
		opt(queryParams)
	}

	httpResp, err := c.call(http.MethodGet, "messages", queryParams, nil)
	if err != nil {
		return nil, err
	}
	var resp struct{ Items []*MessageSummary }
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	return resp.Items, json.Unmarshal(body, &resp)
}

// DeleteMessages permanently deletes all messages held by the specified server.
func (c *Client) DeleteMessages() error {
	_, err := c.call(http.MethodDelete, "messages", map[string]interface{}{
		"server": c.serverID,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

// TODO: implement search messages
