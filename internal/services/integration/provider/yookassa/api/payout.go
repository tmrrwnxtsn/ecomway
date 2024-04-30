package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

const payoutsEndpoint = "/payouts"

type payoutAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type createPayoutRequest struct {
	Amount          payoutAmount `json:"amount"`
	PaymentMethodID string       `json:"payment_method_id"`
	Description     string       `json:"description"`
}

type payoutCancellation struct {
	Party  string `json:"party"`
	Reason string `json:"reason"`
}

type createPayoutResponse struct {
	ID           string             `json:"id"`
	Status       string             `json:"status"`
	CapturedAt   string             `json:"captured_at"`
	Cancellation payoutCancellation `json:"cancellation_details"`
}

func (c *Client) CreatePayout(ctx context.Context, request data.CreatePayoutRequest) (data.CreatePayoutResponse, error) {
	var response data.CreatePayoutResponse

	req := createPayoutRequest{
		Amount: payoutAmount{
			Value:    convert.FloatWithoutTrailingZeroes(request.Amount.Value, amountPrecision),
			Currency: request.Amount.Currency,
		},
		PaymentMethodID: request.PaymentMethodID,
		Description:     request.Description,
	}

	reqPayload, err := json.Marshal(req)
	if err != nil {
		err = fmt.Errorf("marshalling create payout request payload: %w", err)
		return response, err
	}

	reqURL, err := url.JoinPath(c.baseURL, payoutsEndpoint)
	if err != nil {
		return response, fmt.Errorf("failed to join base URL with create payout request endpoint: %w", err)
	}

	reqBody := bytes.NewBuffer(reqPayload)
	reqMethod := http.MethodPost

	httpRequest, err := http.NewRequestWithContext(ctx, reqMethod, reqURL, reqBody)
	if err != nil {
		return response, fmt.Errorf("creating HTTP request with context: %w", err)
	}

	c.setPayoutRequiredHeaders(httpRequest)

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return response, fmt.Errorf("making HTTP request: %w", err)
	}
	defer closeResponseBody(httpResponse)

	respPayload, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		err = errorFromUnresolvedResponse(reqMethod, reqURL, httpResponse.StatusCode, nil)
		return response, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		var resp errorResponse
		if err = json.Unmarshal(respPayload, &resp); err != nil {
			err = errorFromUnresolvedResponse(reqMethod, reqURL, httpResponse.StatusCode, respPayload)
			return response, err
		}

		err = errorResponseData(resp)
		return response, err
	}

	var resp createPayoutResponse
	if err = json.Unmarshal(respPayload, &resp); err != nil {
		err = errorFromUnresolvedResponse(reqMethod, reqURL, httpResponse.StatusCode, respPayload)
		return response, err
	}

	if resp.CapturedAt != "" {
		response.CapturedAt, err = time.Parse(datetimeLayout, resp.CapturedAt)
		if err != nil {
			return response, fmt.Errorf("parse captured_at %q with layout %q: %w", resp.CapturedAt, datetimeLayout, err)
		}
	}

	response.ID = resp.ID
	response.Status = resp.Status
	response.Cancellation = data.Cancellation(resp.Cancellation)

	return response, nil
}
