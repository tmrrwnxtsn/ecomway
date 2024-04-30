package api

import (
	"net/http"
	"time"

	httpclient "github.com/tmrrwnxtsn/ecomway/internal/pkg/http"
)

const requestTimeout = 30 * time.Second

const amountPrecision = 2

const datetimeLayout = "2006-01-02T15:04:05Z"

type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient        HTTPClient
	baseURL           string
	shopID            string
	agentID           string
	paymentsSecretKey string
	payoutsSecretKey  string
}

type ClientOptions struct {
	BaseURL           string
	ShopID            string
	AgentID           string
	PaymentsSecretKey string
	PayoutsSecretKey  string
}

func NewClient(opts ClientOptions) *Client {
	return &Client{
		httpClient:        httpclient.NewClient(requestTimeout, nil),
		baseURL:           opts.BaseURL,
		shopID:            opts.ShopID,
		agentID:           opts.AgentID,
		paymentsSecretKey: opts.PaymentsSecretKey,
		payoutsSecretKey:  opts.PayoutsSecretKey,
	}
}
