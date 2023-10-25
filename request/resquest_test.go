package request

import (
	"fmt"
	"testing"
)

func TestRequest_Get(t *testing.T) {
	resp, _ := Get("https://httpbin.org/get")
	fmt.Println(resp.String())
}
