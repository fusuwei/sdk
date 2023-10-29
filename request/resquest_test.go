package request

import (
	"fmt"
	"testing"
)

func TestRequest_Get(t *testing.T) {
	resp, _ := Get("https://httpbin.org/get")
	fmt.Println(resp.String())
}

func TestPost(t *testing.T) {
	body := map[string]string{
		"1": "11",
	}
	resp, _ := Post("https://httpbin.org/post", nil, body)
	fmt.Println(resp.String())
}
