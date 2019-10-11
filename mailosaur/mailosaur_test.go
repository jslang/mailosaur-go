package mailosaur_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/brianvoe/gofakeit"
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

// NewTestHTTPServer starts an http server that can be used to handle an http request, returns the started service and
// a pointer to an http request that will be used to store the incoming request. The provided TestServiceResponse
func NewTestHTTPServer(resp *TestResponse) (*httptest.Server, *http.Request) {
	var recvReq http.Request
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recvReq = *r
		for key, value := range resp.Headers {
			w.Header().Add(key, value)
		}
		w.WriteHeader(resp.StatusCode)

		if resp.Body == nil {
			return
		}
		if _, err := w.Write(resp.Body); err != nil {
			panic(err)
		}
	}))
	return s, &recvReq
}
