package httpy

import (
	"context"
	"net/http"
	"strings"
	"time"
)

func init() {
	DefaultClient = &standardClient{
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// DefaultClient is the default implementation of Client.
var DefaultClient Client

// Client is a high-level interface used to make HTTP requests.
type Client interface {
	// SetBaseURL sets the base URL for the requests.
	SetBaseURL(url string)

	// SetTimeout sets the timeout for the HTTP requests.
	SetTimeout(d time.Duration)

	// Do sends the HTTP request to the target URL. If a base url is specified
	// it will be prepended the the request's url, however, if the request's
	// url is fully-qualified url, the base url will be ignored.
	Do(ctx context.Context, req *Request) (*Response, error)
}

type standardClient struct {
	baseURL    string
	httpClient *http.Client
}

func (c *standardClient) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *standardClient) SetTimeout(d time.Duration) {
	c.httpClient.Timeout = d
}

func (c *standardClient) Do(ctx context.Context, req *Request) (*Response, error) {
	target := getTargetURL(c.baseURL, req.url)
	r, err := http.NewRequestWithContext(ctx, req.method, target, req.body)
	if err != nil {
		return nil, err
	}
	addRequestHeaders(r, req.headers)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	res := Response(*resp)
	return &res, nil
}

func getTargetURL(base, url string) string {
	if base == "" || strings.HasPrefix(url, "http") {
		return url
	}
	lc := base[len(base)-1]
	if lc == '/' {
		base = base[:len(base)-1]
	}
	if url[0] == '/' {
		url = url[1:]
	}
	return base + "/" + url
}

func addRequestHeaders(r *http.Request, headers map[string][]string) {
	for name, values := range headers {
		for _, value := range values {
			r.Header.Add(name, value)
		}
	}
}
