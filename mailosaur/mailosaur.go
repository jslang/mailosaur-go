// mailosaur is a golang client implementation of the mailosaur test email service API
package mailosaur

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// ServiceURL provides the default service url for the mailosaur API
	ServiceURL = "https://mailosaur.com/api"
)

// Client provides impelmentations of the mailosaur API
type Client struct {
	serverID   string
	apiKey     string
	serviceURL string
	http       http.Client
}

// clientOption is an option function that configures the mailosaur client
type clientOption func(*Client)

// SetServiceURL overrides the default service url the mailosaur client is configured to use.
func SetServiceURL(serviceURL string) clientOption {
	return func(c *Client) {
		c.serviceURL = strings.TrimSuffix(serviceURL, "/")
	}
}

// NewClient creates, configures, and returns a new mailosaur Client
func NewClient(apiKey string, serverID string, options ...clientOption) *Client {
	c := &Client{
		apiKey:     apiKey,
		serverID:   serverID,
		serviceURL: ServiceURL,
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

// call constructs a request to the mailosaur API, applying necessary authorization and request headers to make a
// successful API call.
func (c *Client) call(method string, path string, queryParams map[string]interface{}, data interface{}) (*http.Response, error) {
	req, err := http.NewRequest(method, c.serviceURL+"/"+path, nil)
	if err != nil {
		return nil, err
	}

	setAuthorization(req, c.apiKey)
	setQueryParams(req, queryParams)
	if err := setJSONData(req, data); err != nil {
		return nil, err
	}

	return c.http.Do(req)
}

// setAuthorization provides authorization headers used by the mailosaur API. The API requires HTTP basic auth via a
// generated API key provided as the username.
func setAuthorization(req *http.Request, apiKey string) {
	req.SetBasicAuth(apiKey, "")
}

// setQueryParams sets the query string for a given request based on the provided parameters.
func setQueryParams(req *http.Request, params map[string]interface{}) {
	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = query.Encode()
}

// setJSONData sets the JSON encoded body for a given request based on the provided parameters.
func setJSONData(req *http.Request, data interface{}) error {
	if data == nil {
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	return nil
}
