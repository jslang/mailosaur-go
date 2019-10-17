package mailosaur_test

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/jslang/mailosaur-go/mailosaur"
	"github.com/stretchr/testify/require"
)

func init() {
	gofakeit.Seed(0)
}

// LoadTestData reads testdata files and returns their contents. Causes test failure on error.
func LoadTestData(t *testing.T, name string) []byte {
	b, err := ioutil.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return b
}

// RandomMessageID generates a random message id that can be used for testing
func RandomMessageID() string {
	return gofakeit.UUID()
}

// RandomAPIKey generates a random api key that can be used as a mailosaur API key for testing
func RandomAPIKey() string {
	return gofakeit.UUID()
}

// RandomServerID generates a random server id that can be used as a mailosaur server id for testing
func RandomServerID() string {
	return gofakeit.UUID()
}

// TestResponse is used with NewTestHTTPServer to define the response the server should return for requests made to it.
type TestResponse struct {
	Body       []byte
	StatusCode int
	Headers    map[string]string
}

type ReceivedRequest struct {
	URL     url.URL
	Headers http.Header
	Body    []byte
	Method  string
}

// NewTestHTTPServer starts an http server that can be used to handle an http request, returns the started service and
// a pointer to an http request that will be used to store the incoming request. The provided TestServiceResponse
func NewTestHTTPServer(t *testing.T, resp *TestResponse) (*httptest.Server, *ReceivedRequest) {
	var (
		recvReq ReceivedRequest
		err     error
	)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recvReq.Headers = r.Header
		recvReq.Method = r.Method
		recvReq.URL = *r.URL
		recvReq.Body, err = ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		for key, value := range resp.Headers {
			w.Header().Add(key, value)
		}
		w.WriteHeader(resp.StatusCode)

		if resp.Body == nil {
			return
		}
		_, err = w.Write(resp.Body)
		require.NoError(t, err, "failed to write response body while handling test request")
	}))
	return s, &recvReq
}

func TestCall(t *testing.T) {
	type TestSetup struct {
		apiKey   string
		serverID string
		recvReq  *ReceivedRequest
		client   *mailosaur.Client
	}
	setup := func(t *testing.T) *TestSetup {
		s, recvReq := NewTestHTTPServer(t, &TestResponse{StatusCode: 200})
		apiKey := RandomAPIKey()
		serverID := RandomServerID()
		return &TestSetup{
			apiKey:   apiKey,
			serverID: serverID,
			recvReq:  recvReq,
			client:   mailosaur.NewClient(apiKey, serverID, mailosaur.SetServiceURL(s.URL)),
		}
	}

	t.Run("sends authorization", func(t *testing.T) {
		ts := setup(t)
		_, err := ts.client.Call("GET", "path", nil, nil)
		require.NoError(t, err)

		auth := base64.StdEncoding.EncodeToString([]byte(ts.apiKey + ":"))
		require.Equal(t, "Basic "+auth, ts.recvReq.Headers.Get("Authorization"))
	})

	t.Run("uses requested http method", func(t *testing.T) {
		ts := setup(t)
		_, err := ts.client.Call(http.MethodHead, "path", nil, nil)
		require.NoError(t, err)
		require.Equal(t, http.MethodHead, ts.recvReq.Method)
	})

	t.Run("uses requested path", func(t *testing.T) {
		ts := setup(t)
		_, err := ts.client.Call(http.MethodHead, "thisismypath", nil, nil)
		require.NoError(t, err)
		require.Equal(t, "/thisismypath", ts.recvReq.URL.Path)

	})

	t.Run("sends json encoded body if available", func(t *testing.T) {
		ts := setup(t)
		msgID := RandomMessageID()
		_, err := ts.client.Call(http.MethodPost, "", nil, struct {
			MsgID string
		}{msgID})
		require.NoError(t, err)
		require.JSONEq(t, `{"MsgID": "`+msgID+`"}`, string(ts.recvReq.Body))
	})

	t.Run("sends query params if available", func(t *testing.T) {
		ts := setup(t)
		msgID := RandomMessageID()
		_, err := ts.client.Call(http.MethodPost, "", map[string]interface{}{
			"msgID": msgID,
		}, nil)
		require.NoError(t, err)
		require.Equal(t, "msgID="+msgID, ts.recvReq.URL.Query().Encode())
	})
}

func TestGenerateEmail(t *testing.T) {
	serverID := RandomServerID()
	c := mailosaur.NewClient(RandomAPIKey(), serverID)

	email := c.GenerateEmail()
	parts := strings.SplitN(email, ".", 2)
	require.Len(t, parts, 2)
	require.Len(t, parts[0], 10)
	require.Equal(t, fmt.Sprintf("%s@%s", serverID, mailosaur.SMTPHost), parts[1])
}
