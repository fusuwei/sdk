package request

type retryOption struct {
	MaxRetries int
	RetryHooks []RetryHookFunc
}

type RetryHookFunc func(resp *Response, err error)
