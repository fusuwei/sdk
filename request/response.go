package request

import (
	"net/http"
)

type Response struct {
	*http.Response
	Err     error
	Request *Request
	body    []byte
}
