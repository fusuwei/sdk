package request

import (
	"context"
	"io"
	"net/http"
	"time"
)

type RoundTripper interface {
	RoundTrip(context.Context, *Request) (*Response, error)
}

type Client struct {
	PathParams            map[string]string
	AllowGetMethodPayload bool

	httpClient *http.Client
}

func newClient() *Client {
	t := &http.Transport{
		DisableKeepAlives: false,                          //关闭连接复用，因为后台连接过多最后会造成端口耗尽
		MaxIdleConns:      -1,                             //最大空闲连接数量
		IdleConnTimeout:   time.Duration(3 * time.Second), //空闲连接超时时间
	}

	httpClient := &http.Client{
		Transport: t,
		Timeout:   2 * time.Minute,
	}

	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) roundTrip(ctx context.Context, r *Request) (resp *Response, err error) {
	resp = &Response{Request: r}

	host := r.getHeader("Host")
	if host == "" {
		host = r.URL.Host
	}

	if r.Timeout != 0 {
		c.httpClient.Timeout = r.Timeout
	}

	var reqBody io.ReadCloser
	if r.GetBody != nil {
		reqBody, err = r.GetBody()
		if err != nil {
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
		ContentLength: int64(len(r.Body)),
		Body:          reqBody,
		GetBody:       r.GetBody,
		Close:         r.close,
	}

	for _, cookie := range r.Cookies {
		req.AddCookie(cookie)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	r.RawRequest = req

	var httpResponse *http.Response
	httpResponse, err = c.httpClient.Do(r.RawRequest)
	resp.Response = httpResponse
	return
}

func (c *Client) isPayloadForbid(m string) bool {
	return (m == http.MethodGet && !c.AllowGetMethodPayload) || m == http.MethodHead || m == http.MethodOptions
}
