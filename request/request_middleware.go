package request

import "net/http"

type (
	// RequestMiddleware type is for request middleware, called before a request is sent
	RequestMiddleware func(client *Client, req *Request) error

	// ResponseMiddleware type is for response middleware, called after a response has been received
	ResponseMiddleware func(client *Client, resp *Response) error
)

func parseRequestBody(c *Client, r *Request) (err error) {
	if r.GetBody != nil {
		return
	}

	if c.isPayloadForbid(r.Method) {
		r.unMarshalBody = nil
		r.Body = nil
		r.GetBody = nil
		return
	}
	//if r.isMultiPart {
	//	return handleMultiPart(c, r)
	//}
	//
	//// handle form data
	//if len(c.FormData) > 0 {
	//	r.SetFormDataFromValues(c.FormData)
	//}
	//if len(r.FormData) > 0 {
	//	handleFormData(r)
	//	return
	//}

	if r.unMarshalBody != nil {
		err = r.handleMarshalBody(r.Body)
		if err != nil {
			return
		}
	}

	if r.Body == nil {
		return
	}
	// body is in-memory []byte, so we can guess content type
	if r.getHeader(ContentType) != "" {
		return
	}
	r.SetContentType(http.DetectContentType(r.Body))
	return
}
