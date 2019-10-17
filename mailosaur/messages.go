package mailosaur

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// messageListOption configures options for listing message lists, for setting page, items per page, time filters, etc.
type messageListOption func(map[string]interface{})

// SetPage sets the page of results to request
func SetPage(page int) messageListOption {
	return func(data map[string]interface{}) {
		data["page"] = page
	}
}

// SetItemsPerPage sets the number of items to display per result page
func SetItemsPerPage(itemsPerPage int) messageListOption {
	return func(data map[string]interface{}) {
		data["itemsPerPage"] = itemsPerPage
	}
}

// SetReceivedAfter sets the time to filter messages by
func SetReceivedAfter(receivedAfter time.Time) messageListOption {
	return func(data map[string]interface{}) {
		data["receivedAfter"] = receivedAfter.Format(time.RFC3339)
	}
}

func applyMessageListOptions(data map[string]interface{}, options []messageListOption) {
	for _, opt := range options {
		opt(data)
	}
}

// GetMessage retrieves the detail for a single email message.
func (c *Client) GetMessage(messageID string) (*Message, error) {
	httpResp, err := c.Call(http.MethodGet, "messages/"+messageID, nil, nil)
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
	_, err := c.Call(http.MethodDelete, "messages/"+messageID, nil, nil)
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

// ListMessages returns a list of your messages in summary form.
func (c *Client) ListMessages(options ...messageListOption) ([]*MessageSummary, error) {
	queryParams := map[string]interface{}{
		"server": c.serverID,
	}
	applyMessageListOptions(queryParams, options)

	httpResp, err := c.Call(http.MethodGet, "messages", queryParams, nil)
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
	_, err := c.Call(http.MethodDelete, "messages", map[string]interface{}{
		"server": c.serverID,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

// SearchMessagesLookup defines the search parameters for a SearchMessages call.
type SearchMessagesLookup struct {
	SentTo  string `json:"sentTo,omitempty"`
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body,omitempty"`
}

// SearchMessages returns a list of message summaries matching the specified search criteria.
func (c *Client) SearchMessages(lookup *SearchMessagesLookup, options ...messageListOption) ([]*MessageSummary, error) {
	queryParams := map[string]interface{}{
		"server": c.serverID,
	}
	applyMessageListOptions(queryParams, options)

	httpResp, err := c.Call(http.MethodPost, "messages/search", queryParams, lookup)
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
