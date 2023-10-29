package request

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
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

func Test_Url(t *testing.T) {
	resp, err := New().
		Request(http.MethodGet, "http://fund.eastmoney.com/pingzhongdata/003885.js").
		SetQuery(map[string]string{"v": fmt.Sprintf("%d", time.Now().UnixMicro())}).
		SetHeaders(map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"}).
		Do(context.Background())
	fmt.Println(err)
	fmt.Println(resp.StatusCode)
	a := resp.String()
	t.Log(a)
	fmt.Println(a)
}
