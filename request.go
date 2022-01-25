package httpy

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
)

const (
	ContentTypeHeader   = "Content-Type"
	AuthorizationHeader = "Authorization"
)

// Request contains a set of data require to make a request. As well
// as providing a fluent-api for building HTTP requests.
type Request struct {
	body    io.Reader
	url     string
	method  string
	headers map[string][]string
}

// NewRequest returns a new Request with the given url and method.
func NewRequest(url, method string) *Request {
	return &Request{
		url:     url,
		method:  method,
		headers: map[string][]string{},
	}
}

// URL returns the Request's configured URL value.
func (r *Request) URL() string {
	return r.url
}

// Method returns the Request's configured HTTP method value.
func (r *Request) Method() string {
	return r.method
}

// Headers returns the Request's header values.
func (r *Request) Headers() map[string][]string {
	return r.headers
}

// Body returns the Request's body value. Unless configured beforehand
// this value will be nil.
func (r *Request) Body() io.Reader {
	return r.body
}

// WithJSON sets the body of the Request to the given data encoded as JSON -
// also sets the Content-Type header to "application/json".
func (r *Request) WithJSON(data interface{}) *Request {
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(data)
	return r.SetBody(buf).SetContentType("application/json")
}

// WithPlainText sets the body of the Request to the string data given, then
// sets the Content-Type header to "text/plain".
func (r *Request) WithPlainText(txt string) *Request {
	return r.SetBody(strings.NewReader(txt)).SetContentType("text/plain")
}

// WithForm sets the body of the Request to the encodes url values provided.
// The content type of the request is also set to "application/x-www-form-urlencoded".
func (r *Request) WithForm(values url.Values) *Request {
	rdr := strings.NewReader(values.Encode())
	return r.SetBody(rdr).SetContentType("application/x-www-form-urlencoded")
}

// SetBody is used to set the body of the Request. If the body has
// already been set, it will be overwritten.
func (r *Request) SetBody(rdr io.Reader) *Request {
	r.body = rdr
	return r
}

// SetContentType sets the Content-Type header of the request.
func (r *Request) SetContentType(value string) *Request {
	return r.SetHeader(ContentTypeHeader, value)
}

func (r *Request) WithBearer(token string) *Request {
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = token[7:]
	}
	return r.SetHeader(AuthorizationHeader, fmt.Sprintf("Bearer %s", token))
}

func (r *Request) WithBasicAuth(username, password string) *Request {
	auth := fmt.Sprintf("%s:%s", username, password)
	value := base64.StdEncoding.EncodeToString([]byte(auth))
	return r.SetHeader(AuthorizationHeader, fmt.Sprintf("Basic %s", value))
}

// SetHeader sets a request header with the given name and specified values.
func (r *Request) SetHeader(name string, values ...string) *Request {
	if r.headers == nil {
		r.headers = map[string][]string{}
	}
	if values == nil {
		delete(r.headers, name)
	} else {
		r.headers[name] = values
	}
	return r
}

// RemoveHeader is used to remove a header from the Request by the specified name.
func (r *Request) RemoveHeader(name string) *Request {
	return r.SetHeader(name)
}

// Do sends the HTTP Request using the DefaultClient, unless provided with another.
func (r *Request) Do(ctx context.Context, client ...Client) (*Response, error) {
	c := DefaultClient
	if len(client) > 0 {
		c = client[0]
	}
	return c.Do(ctx, r)
}
