package request

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

type Request struct {
	Headers    http.Header    // 请求头
	RawRequest *http.Request  // 请求对象
	RowUrl     string         // 请求地址
	URL        *url.URL       // 请求地址结构体
	Method     string         // 请求方式
	Body       []byte         // 请求请求主体
	GetBody    GetContentFunc // 获取请求主体方法
	Cookies    []*http.Cookie // cookie
	Timeout    time.Duration  // 超时时间
	// 重试
	retryOption *retryOption
	client      *Client
	close       bool
	// 请求，支持自定义封装
	roundTrip RoundTripper

	unReplayableBody io.ReadCloser
	unMarshalBody    interface{} // 没有解析的body, 对象，结构体

	beforeRequest []RequestMiddleware
	afterResponse []ResponseMiddleware
}

type GetContentFunc func() (io.ReadCloser, error)

func New() *Request {
	beforeRequest := []RequestMiddleware{
		parseRequestBody,
	}

	return &Request{
		Headers:          http.Header{},
		Cookies:          make([]*http.Cookie, 0),
		retryOption:      nil,
		client:           nil,
		close:            false,
		roundTrip:        nil,
		unReplayableBody: nil,
		unMarshalBody:    nil,
		beforeRequest:    beforeRequest,
		afterResponse:    make([]ResponseMiddleware, 0),
	}
}

func (r *Request) Do(ctx context.Context) (resp *Response, err error) {
	if r.client == nil {
		r.client = newClient()
	}

	resp, err = r.do(ctx)

	return
}

func (r *Request) do(ctx context.Context) (resp *Response, err error) {
	if r.Headers == nil {
		r.Headers = make(http.Header)
	}
	for _, f := range r.beforeRequest {
		if err = f(r.client, r); err != nil {
			return
		}
	}
	resp, err = r.client.roundTrip(ctx, r)
	for _, f := range r.afterResponse {
		if err = f(r.client, resp); err != nil {
			return
		}
	}
	return
}

func (r *Request) Request(method, rowUrl string) (resp *Response, err error) {
	r.URL, err = url.Parse(rowUrl)
	if err != nil {
		return
	}
	r.Method = method
	r.RowUrl = rowUrl
	resp, err = r.Do(context.Background())
	return
}

func (r *Request) SetBody(body interface{}) *Request {
	if body == nil {
		return r
	}
	switch b := body.(type) {
	case io.ReadCloser:
		r.unReplayableBody = b
		r.GetBody = func() (io.ReadCloser, error) {
			return r.unReplayableBody, nil
		}
	case io.Reader:
		r.unReplayableBody = io.NopCloser(b)
		r.GetBody = func() (io.ReadCloser, error) {
			return r.unReplayableBody, nil
		}
	case []byte:
		r.SetBodyBytes(b)
	case string:
		r.SetBodyString(b)
	case func() (io.ReadCloser, error):
		r.GetBody = b
	case GetContentFunc:
		r.GetBody = b
	default:
		t := reflect.TypeOf(body)
		switch t.Kind() {
		case reflect.Ptr, reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
			r.unMarshalBody = body
		default:
			r.SetBodyString(fmt.Sprint(body))
		}
	}
	return r
}

func (r *Request) handleMarshalBody(body interface{}) error {
	ct := ""
	if r.Headers != nil {
		ct = r.Headers.Get(ContentType)
	}
	if ct != "" {
		if IsXMLType(ct) {
			body, err := xml.Marshal(body)
			if err != nil {
				return err
			}
			r.SetBodyBytes(body)
		} else {
			body, err := json.Marshal(body)
			if err != nil {
				return err
			}
			r.SetBodyBytes(body)
		}
		return nil
	}
	content, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r.SetBodyJsonBytes(content)
	return nil
}

func (r *Request) SetBodyBytes(body []byte) *Request {
	r.Body = body
	r.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}
	return r
}

func (r *Request) SetBodyString(body string) *Request {
	return r.SetBodyBytes([]byte(body))
}

func (r *Request) SetBodyJsonBytes(body []byte) *Request {
	r.SetContentType(JsonContentType)
	return r.SetBodyBytes(body)
}

func (r *Request) SetContentType(contentType string) *Request {
	return r.SetHeader(ContentType, contentType)
}

func (r *Request) SetFormData(data map[string]string) *Request {
	if len(data) == 0 {
		return r
	}
	formData := url.Values{}
	for k, v := range data {
		formData.Add(k, v)
	}
	r.SetContentType(FormContentType)
	r.SetBodyBytes([]byte(formData.Encode()))
	return r
}

func (r *Request) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		r.SetHeader(k, v)
	}
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.Headers == nil {
		r.Headers = make(http.Header)
	}
	r.Headers.Set(key, value)
	return r
}

func (r *Request) getHeader(key string) string {
	if r.Headers == nil {
		return ""
	}
	return r.Headers.Get(key)
}

func Get(url string) (resp *Response, err error) {
	return New().Request(http.MethodGet, url)
}

func Post(url string, body interface{}, form map[string]string) (resp *Response, err error) {
	req := New()
	if body != nil {
		req.SetBody(body).SetContentType(JsonContentType)
	}
	if form != nil {
		req.SetFormData(form)
	}
	return req.Request(http.MethodPost, url)
}
