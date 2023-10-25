package request

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	*http.Response
	Request *Request
	body    []byte
}

func (r *Response) UnmarshalJson(v interface{}) error {
	b, err := r.ToBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (r *Response) ToBytes() (body []byte, err error) {
	if r.body != nil {
		return r.body, nil
	}
	if r.Response == nil || r.Response.Body == nil {
		return []byte{}, nil
	}
	defer func() {
		r.Body.Close()
		if err != nil {
			return
		}
		r.body = body
	}()
	body, err = io.ReadAll(r.Body)
	return
}

func (r *Response) String() string {
	body, err := r.ToBytes()
	if err != nil {
		return ""
	}
	return string(body)
}
