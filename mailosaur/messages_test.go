package mailosaur_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/jslang/mailosaur-go/mailosaur"
	"github.com/stretchr/testify/require"
)

func TestGetMessage(t *testing.T) {
	type testSetup struct {
		apiKey   string
		serverID string
		recvReq  *http.Request
		client   *mailosaur.Client
	}

	setup := func(t *testing.T, resp *TestResponse) *testSetup {
		t.Parallel()
		s, recvReq := NewTestHTTPServer(resp)
		apiKey := RandomAPIKey()
		serverID := RandomServerID()
		return &testSetup{
			apiKey:   apiKey,
			serverID: serverID,
			recvReq:  recvReq,
			client:   mailosaur.NewClient(apiKey, serverID, mailosaur.SetServiceURL(s.URL)),
		}
	}

	t.Run("calls get message endpoint", func(t *testing.T) {
		ts := setup(t, &TestResponse{
			Body:       LoadTestData(t, "get_message_success.json"),
			StatusCode: http.StatusOK,
		})

		msgID := RandomMessageID()
		_, err := ts.client.GetMessage(msgID)
		require.NoError(t, err)
		require.Equal(t, "/messages/"+msgID, ts.recvReq.URL.Path)
		require.Equal(t, http.MethodGet, ts.recvReq.Method)
	})
}
func TestDeleteMessage(t *testing.T) {
	type testSetup struct {
		apiKey   string
		serverID string
		recvReq  *http.Request
		client   *mailosaur.Client
	}

	setup := func(t *testing.T, resp *TestResponse) *testSetup {
		t.Parallel()
		s, recvReq := NewTestHTTPServer(resp)
		apiKey := RandomAPIKey()
		serverID := RandomServerID()
		return &testSetup{
			apiKey:   apiKey,
			serverID: serverID,
			recvReq:  recvReq,
			client:   mailosaur.NewClient(apiKey, serverID, mailosaur.SetServiceURL(s.URL)),
		}
	}

	t.Run("calls delete message endpoint", func(t *testing.T) {
		ts := setup(t, &TestResponse{
			StatusCode: http.StatusNoContent,
		})

		msgID := RandomMessageID()
		err := ts.client.DeleteMessage(msgID)
		require.NoError(t, err)
		require.Equal(t, "/messages/"+msgID, ts.recvReq.URL.Path)
		require.Equal(t, http.MethodDelete, ts.recvReq.Method)
	})
}

func TestListMessages(t *testing.T) {
	type testSetup struct {
		apiKey   string
		serverID string
		recvReq  *http.Request
		client   *mailosaur.Client
	}

	setup := func(t *testing.T, resp *TestResponse) *testSetup {
		t.Parallel()
		s, recvReq := NewTestHTTPServer(resp)
		apiKey := RandomAPIKey()
		serverID := RandomServerID()
		return &testSetup{
			apiKey:   apiKey,
			serverID: serverID,
			recvReq:  recvReq,
			client:   mailosaur.NewClient(apiKey, serverID, mailosaur.SetServiceURL(s.URL)),
		}
	}

	t.Run("calls list messages endpoint", func(t *testing.T) {
		ts := setup(t, &TestResponse{
			Body:       LoadTestData(t, "list_messages_success.json"),
			StatusCode: http.StatusOK,
		})

		_, err := ts.client.ListMessages()
		require.NoError(t, err)
		require.Equal(t, "/messages", ts.recvReq.URL.Path)
		require.Equal(t, http.MethodGet, ts.recvReq.Method)
	})

	t.Run("uses configured server id", func(t *testing.T) {
		ts := setup(t, &TestResponse{
			Body:       LoadTestData(t, "list_messages_success.json"),
			StatusCode: http.StatusOK,
		})

		_, err := ts.client.ListMessages()
		require.NoError(t, err)
		require.Equal(t, ts.serverID, ts.recvReq.URL.Query().Get("server"))
	})

	t.Run("uses provided page", func(t *testing.T) {
		ts := setup(t, &TestResponse{
			Body:       LoadTestData(t, "list_messages_success.json"),
			StatusCode: http.StatusOK,
		})

		_, err := ts.client.ListMessages(mailosaur.SetListMessagesPage(10))
		require.NoError(t, err)
		require.Equal(t, "10", ts.recvReq.URL.Query().Get("page"))
	})

	t.Run("uses provided items per page", func(t *testing.T) {
		ts := setup(t, &TestResponse{
			Body:       LoadTestData(t, "list_messages_success.json"),
			StatusCode: http.StatusOK,
		})

		_, err := ts.client.ListMessages(mailosaur.SetListMessagesItemsPerPage(100))
		require.NoError(t, err)
		require.Equal(t, "100", ts.recvReq.URL.Query().Get("itemsPerPage"))
	})

	t.Run("uses provided receivedAfter", func(t *testing.T) {
		ts := setup(t, &TestResponse{
			Body:       LoadTestData(t, "list_messages_success.json"),
			StatusCode: http.StatusOK,
		})

		receivedAfter := time.Now()
		_, err := ts.client.ListMessages(mailosaur.SetListMessagesReceivedAfter(receivedAfter))
		require.NoError(t, err)
		require.Equal(t, receivedAfter.Format(time.RFC3339), ts.recvReq.URL.Query().Get("receivedAfter"))
	})
}
func TestDeleteMessages(t *testing.T) {
	type testSetup struct {
		apiKey   string
		serverID string
		recvReq  *http.Request
		client   *mailosaur.Client
	}

	setup := func(t *testing.T, resp *TestResponse) *testSetup {
		t.Parallel()
		s, recvReq := NewTestHTTPServer(resp)
		apiKey := RandomAPIKey()
		serverID := RandomServerID()
		return &testSetup{
			apiKey:   apiKey,
			serverID: serverID,
			recvReq:  recvReq,
			client:   mailosaur.NewClient(apiKey, serverID, mailosaur.SetServiceURL(s.URL)),
		}
	}

	t.Run("calls delete message endpoint", func(t *testing.T) {
		ts := setup(t, &TestResponse{
			StatusCode: http.StatusNoContent,
		})

		err := ts.client.DeleteMessages()
		require.NoError(t, err)
		require.Equal(t, "/messages", ts.recvReq.URL.Path)
		require.Equal(t, http.MethodDelete, ts.recvReq.Method)
	})

	t.Run("uses configured server id", func(t *testing.T) {
		ts := setup(t, &TestResponse{
			StatusCode: http.StatusNoContent,
		})

		err := ts.client.DeleteMessages()
		require.NoError(t, err)
		require.Equal(t, ts.serverID, ts.recvReq.URL.Query().Get("server"))
	})
}
