package api

import (
	"net/http"
	"time"

	httpclient "github.com/tmrrwnxtsn/ecomway/internal/pkg/http"
)

const requestTimeout = 30 * time.Second

const amountPrecision = 2

type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HTTPClient
	baseURL    string
	shopID     string
	secretKey  string
}

type ClientOptions struct {
	BaseURL   string
	ShopID    string
	SecretKey string
}

func NewClient(opts ClientOptions) *Client {
	return &Client{
		httpClient: httpclient.NewClient(requestTimeout, nil),
		baseURL:    opts.BaseURL,
		shopID:     opts.ShopID,
		secretKey:  opts.SecretKey,
	}
}
