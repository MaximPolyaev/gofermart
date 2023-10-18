package accrual

import "net/http"

type HTTPClient struct {
	client http.Client
}

func New() *HTTPClient {
	return &HTTPClient{client: http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}}
}
