package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

const createPaymentEndpoint = "/payments"

type paymentAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type paymentMethodData struct {
	Type string `json:"type"`
}

type paymentConfirmation struct {
	Type            string `json:"type"`
	Locale          string `yaml:"locale,omitempty"`
	ReturnURL       string `json:"return_url"`                 // URL для возврата пользователя после оплаты
	ConfirmationURL string `json:"confirmation_url,omitempty"` // URL для перенаправления пользователя на страницу оплаты
}

type createPaymentRequest struct {
	Amount            paymentAmount       `json:"amount"`
	Capture           bool                `yaml:"capture"`
	PaymentMethodData paymentMethodData   `json:"payment_method_data"`
	Confirmation      paymentConfirmation `json:"confirmation"`
	Description       string              `json:"description,omitempty"`
}

type createPaymentResponse struct {
	ID           string              `json:"id"`
	Status       string              `json:"status"`
	Confirmation paymentConfirmation `json:"confirmation"`
}

func (c *Client) CreatePayment(ctx context.Context, request data.CreatePaymentRequest) (data.CreatePaymentResponse, error) {
	var response data.CreatePaymentResponse

	req := createPaymentRequest{
		Amount: paymentAmount{
			Value:    convert.FloatWithoutTrailingZeroes(request.Amount.Amount, amountPrecision),
			Currency: request.Amount.Currency,
		},
		Capture: request.Capture,
		PaymentMethodData: paymentMethodData{
			Type: request.PaymentMethod.Type,
		},
		Confirmation: paymentConfirmation{
			Type:      request.Confirmation.Type,
			Locale:    request.Confirmation.Locale,
			ReturnURL: request.Confirmation.ReturnURL,
		},
		Description: request.Description,
	}

	reqPayload, err := json.Marshal(req)
	if err != nil {
		err = fmt.Errorf("marshalling create payment request payload: %w", err)
		return response, err
	}

	reqURL, err := url.JoinPath(c.baseURL, createPaymentEndpoint)
	if err != nil {
		return response, fmt.Errorf("failed to join base URL with create payment request endpoint: %w", err)
	}

	reqBody := bytes.NewBuffer(reqPayload)
	reqMethod := http.MethodPost

	httpRequest, err := http.NewRequestWithContext(ctx, reqMethod, reqURL, reqBody)
	if err != nil {
		return response, fmt.Errorf("creating HTTP request with context: %w", err)
	}

	c.setRequiredHeaders(httpRequest)

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
		// TODO: обработать ошибочный ответ: https://yookassa.ru/developers/using-api/response-handling/response-format#error
		return response, errors.New("unexpected error")
	}

	var resp createPaymentResponse
	if err = json.Unmarshal(respPayload, &resp); err != nil {
		err = errorFromUnresolvedResponse(reqMethod, reqURL, httpResponse.StatusCode, respPayload)
		return response, err
	}

	response.ID = resp.ID
	response.ConfirmationURL = resp.Confirmation.ConfirmationURL
	response.Status = resp.Status

	return response, nil
}
