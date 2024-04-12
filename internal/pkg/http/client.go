package http

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient(timeout time.Duration, transport http.RoundTripper) *Client {
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	return &Client{
		client: httpClient,
	}
}

func (c *Client) Do(request *http.Request) (*http.Response, error) {
	var body []byte
	if request.Body != nil {
		bodyReaderCopy, err := request.GetBody()
		if err != nil {
			return nil, err
		}
		body, err = io.ReadAll(bodyReaderCopy)
		if err != nil {
			return nil, err
		}
	}

	slog.Info(
		"outgoing request",
		"method", request.Method,
		"url", request.URL.String(),
		"body", string(body),
	)

	response, err := c.client.Do(request)
	if err != nil {
		slog.Error(
			"request error",
			"method", request.Method,
			"url", request.URL.String(),
			"body", string(body),
			"error", err,
		)
		return nil, err
	}

	responseBodyData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	slog.Info(
		"response to outgoing request",
		"status_code", response.StatusCode,
		"url", request.URL.String(),
		"body", string(responseBodyData),
	)

	response.Body = io.NopCloser(bytes.NewReader(responseBodyData))

	return response, nil
}
