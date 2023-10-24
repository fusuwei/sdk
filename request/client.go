package request

import (
	"context"
	"io"
	"net/http"
)

type RoundTripper interface {
	RoundTrip(context.Context, *Request) (*Response, error)
}

type Client struct {
	PathParams            map[string]string
	DebugLog              bool
	AllowGetMethodPayload bool

	beforeRequest []RequestMiddleware
	afterResponse []ResponseMiddleware

	httpClient *http.Client
}

func (c *Client) RoundTrip(ctx context.Context, r *Request) (*Response, error) {
	resp = &Response{Request: r}

	host := r.getHeader("Host")
	if host == "" {
		host = r.URL.Host
	}

	var reqBody io.ReadCloser
	if r.GetBody != nil {
		reqBody, resp.Err = r.GetBody()
		if resp.Err != nil {
			return
		}
	}

	req := &http.Request{
		Method:        r.Method,
		Header:        r.Headers.Clone(),
		URL:           r.URL,
		Host:          host,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: contentLength,
		Body:          reqBody,
		GetBody:       r.GetBody,
		Close:         r.close,
	}

	return nil, nil
}

func (c *Client) isPayloadForbid(m string) bool {
	return (m == http.MethodGet && !c.AllowGetMethodPayload) || m == http.MethodHead || m == http.MethodOptions
}
