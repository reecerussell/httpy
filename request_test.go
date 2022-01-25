package httpy_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	. "github.com/reecerussell/httpy"
	. "github.com/reecerussell/httpy/mock"
)

func TestNewRequest_GivenURLAndMethod_ReturnsRequest(t *testing.T) {
	r := NewRequest("foo", "bar")
	assert.Equal(t, "foo", r.URL())
	assert.Equal(t, "bar", r.Method())
	assert.NotNil(t, r.Headers())
}

func TestRequestWithJSON_GivenData_MarshalsAndAddsToRequest(t *testing.T) {
	data := map[string]string{
		"foo": "bar",
	}
	r := &Request{}
	r.WithJSON(data)

	assert.NotNil(t, r.Body())
	assert.Contains(t, r.Headers()[ContentTypeHeader], "application/json")

	body := map[string]string{}
	_ = json.NewDecoder(r.Body()).Decode(&body)

	assert.Equal(t, "bar", body["foo"])
}

func TestRequestWithPlainText_GivenString_SetsBody(t *testing.T) {
	const text = "hello world"

	r := &Request{}
	r.WithPlainText(text)

	assert.NotNil(t, r.Body())
	assert.Contains(t, r.Headers()[ContentTypeHeader], "text/plain")

	bytes, _ := ioutil.ReadAll(r.Body())
	assert.Equal(t, text, string(bytes))
}

func TestRequestWithForm_GivenValues_SetsBody(t *testing.T) {
	r := &Request{}
	r.WithForm(url.Values{"msg": {"hello"}})

	assert.NotNil(t, r.Body())
	assert.Contains(t, r.Headers()[ContentTypeHeader], "application/x-www-form-urlencoded")

	bytes, _ := ioutil.ReadAll(r.Body())
	values, _ := url.ParseQuery(string(bytes))

	assert.Equal(t, "hello", values.Get("msg"))
}

func TestRequestWithBearer_GivenTokenWithBearerPrefix_AddsAuthHeader(t *testing.T) {
	r := &Request{}
	r.WithBearer("Bearer abc")

	assert.Contains(t, r.Headers()[AuthorizationHeader], "Bearer abc")
}

func TestRequestWithBearer_GivenTokenWithoutBearerPrefix_AddsAuthHeader(t *testing.T) {
	r := &Request{}
	r.WithBearer("abc")

	assert.Contains(t, r.Headers()[AuthorizationHeader], "Bearer abc")
}

func TestRequestWithBasicAuth_GivenUsernameAndPassword_SetsBasicAuthHeader(t *testing.T) {
	r := &Request{}
	r.WithBasicAuth("foo", "bar")

	assert.Contains(t, r.Headers()[AuthorizationHeader], "Basic Zm9vOmJhcg==")
}

func TestRequestRemoveHeader_GivenValidName_RemovesHeader(t *testing.T) {
	r := &Request{}
	r.SetHeader("foo", "bar")
	assert.Equal(t, []string{"bar"}, r.Headers()["foo"])

	r.RemoveHeader("foo")

	assert.Nil(t, r.Headers()["foo"])
}

func TestRequestDo_GivenAClient_MakesRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := NewRequest("http://localhost/api", http.MethodGet)
	testResp := &Response{}

	client := NewMockClient(ctrl)
	client.EXPECT().Do(ctx, req).Return(testResp, nil)

	resp, err := req.Do(ctx, client)
	assert.Equal(t, testResp, resp)
	assert.Nil(t, err)
}
