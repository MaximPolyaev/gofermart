package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/cenkalti/backoff/v4"
	"io"
	"net/http"
)

const requestAttemptsCount = 15

type HTTPClient struct {
	client  *http.Client
	baseURL string
}

func New(baseURL string) *HTTPClient {
	return &HTTPClient{client: &http.Client{}, baseURL: baseURL}
}

func (c *HTTPClient) FetchAccrualOrder(ctx context.Context, number string) (*entities.AccrualOrder, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/api/orders/"+number,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.repeatableDo(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, errors.New("order is not registered")
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("status many reqs")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var accrualOrder entities.AccrualOrder
	err = json.Unmarshal(body, &accrualOrder)
	if err != nil {
		return nil, err
	}

	return &accrualOrder, nil
}

func (c *HTTPClient) repeatableDo(r *http.Request) (*http.Response, error) {
	ticker := backoff.NewTicker(backoff.NewExponentialBackOff())

	var resp *http.Response
	var err error

	attempts := 0
	for range ticker.C {
		attempts++

		resp, err = c.client.Do(r)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusTooManyRequests && requestAttemptsCount > attempts {
			continue
		}

		ticker.Stop()
		break
	}

	return resp, err
}
