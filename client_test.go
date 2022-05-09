package httpy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientDo_GivenValidRequest_MakesRequest(t *testing.T) {
	// Times the request has been received.
	var times int32 = 0

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&times, 1)

		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api", r.RequestURI)
		assert.Equal(t, []string{"bar", "foobar"}, r.Header.Values("foo"))

		// Ensure body is correct
		payload := map[string]string{}
		json.NewDecoder(r.Body).Decode(&payload)
		assert.Equal(t, "hello world", payload["message"])

		data := map[string]string{"message": "greetings"}
		json.NewEncoder(rw).Encode(data)
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	payload := map[string]string{"message": "hello world"}
	req := NewRequest("/api", http.MethodPost).
		SetHeader("foo", "bar", "foobar").
		WithJSON(payload)

	DefaultClient.SetBaseURL(server.URL)
	resp, err := DefaultClient.Do(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(1), times)

	body := map[string]string{}
	err = resp.DecodeJSON(&body)
	assert.Nil(t, err)
	assert.Equal(t, "greetings", body["message"])
}

func TestClientDo_GivenInvalidMethod_FailsToMakeHTTPRequest(t *testing.T) {
	req := NewRequest("http://localhost", "MY REQUEST METHOD")

	resp, err := DefaultClient.Do(context.Background(), req)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}

func TestClientDo_GivenInvalidURL_FailsToSendRequest(t *testing.T) {
	req := NewRequest("not a url", http.MethodGet)

	resp, err := DefaultClient.Do(context.Background(), req)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}

func TestClientSetTimeout_GivenDuration_SetsUnderlyingClientTimeout(t *testing.T) {
	DefaultClient.SetTimeout(time.Second * 10)

	t.Cleanup(func() {
		DefaultClient.(*standardClient).httpClient.Timeout = time.Second * 30
	})

	assert.Equal(t, time.Second*10, DefaultClient.(*standardClient).httpClient.Timeout)
}

func TestGetTargetURL_WhereBaseIsEmpty_ReturnsURL(t *testing.T) {
	const url = "/api"

	res := getTargetURL("", url)
	assert.Equal(t, url, res)
}

func TestGetTargetURL_WhereURLIsFullyQualified_ReturnsURL(t *testing.T) {
	const url = "https://localhost/api"

	res := getTargetURL("http://google.com", url)
	assert.Equal(t, url, res)
}

func TestGetTargetURL_GivenDifferentVariations_HandlesLeadingAndTrailingSlashes(t *testing.T) {
	// No leading or trailing slash
	base := "http://localhost/api"
	url := "values"
	res := getTargetURL(base, url)
	assert.Equal(t, "http://localhost/api/values", res)

	// Leading slash
	base = "http://localhost/api/"
	url = "values"
	res = getTargetURL(base, url)
	assert.Equal(t, "http://localhost/api/values", res)

	// Trailing slash
	base = "http://localhost/api"
	url = "/values"
	res = getTargetURL(base, url)
	assert.Equal(t, "http://localhost/api/values", res)

	// Leading and trailing slashes
	base = "http://localhost/api/"
	url = "/values"
	res = getTargetURL(base, url)
	assert.Equal(t, "http://localhost/api/values", res)
}

func TestNewClient_GivenClient_CreatesWithGivenClient(t *testing.T) {
	hc := &http.Client{}
	c := NewClient(hc)
	assert.Equal(t, hc, c.(*standardClient).httpClient)
}

func TestNewClient_NotGivenClient_CreatesNewClient(t *testing.T) {
	c := NewClient()
	assert.NotNil(t, c.(*standardClient).httpClient)
}
