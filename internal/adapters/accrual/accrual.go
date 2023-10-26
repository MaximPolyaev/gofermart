package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/errors/accrualerrors"
)

const defaultRetryAfter = 60

type HTTPClient struct {
	client  *http.Client
	baseURL string
	log     logger
}

type logger interface {
	Error(args ...interface{})
}

func New(baseURL string, log logger) *HTTPClient {
	return &HTTPClient{client: &http.Client{}, baseURL: baseURL, log: log}
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

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, &accrualerrors.RateLimitError{RetryAfter: 10}
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := c.getRetryAfterByResp(resp)
		return nil, &accrualerrors.RateLimitError{RetryAfter: retryAfter}
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

func (c *HTTPClient) getRetryAfterByResp(r *http.Response) int {
	respRetryAfter := r.Header.Get("Retry-After")
	if respRetryAfter != "" {
		retryAfter, err := strconv.Atoi(respRetryAfter)
		if err == nil {
			return retryAfter
		}

		c.log.Error(err)
	}

	return defaultRetryAfter
}
