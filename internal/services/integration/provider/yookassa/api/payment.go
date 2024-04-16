package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/convert"
	"github.com/tmrrwnxtsn/ecomway/internal/services/integration/provider/yookassa/data"
)

const paymentsEndpoint = "/payments"

const datetimeLayout = "2006-01-02T15:04:05Z"

type paymentAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type paymentMethodCard struct {
	First6      string `json:"first6"`
	Last4       string `json:"last4"`
	ExpiryYear  string `json:"expiry_year"`
	ExpiryMonth string `json:"expiry_month"`
	CardType    string `json:"card_type"`
	IssuerName  string `json:"issuer_name"`
}

type paymentMethod struct {
	Type  string             `json:"type"`
	ID    string             `json:"id,omitempty"`
	Saved bool               `json:"saved,omitempty"`
	Card  *paymentMethodCard `json:"card,omitempty"`
}

type paymentConfirmation struct {
	Type            string `json:"type"`
	Locale          string `json:"locale,omitempty"`
	ReturnURL       string `json:"return_url"`                 // URL для возврата пользователя после оплаты
	ConfirmationURL string `json:"confirmation_url,omitempty"` // URL для перенаправления пользователя на страницу оплаты
}

type createPaymentRequest struct {
	Amount            paymentAmount       `json:"amount"`
	Capture           bool                `json:"capture"`
	SavePaymentMethod bool                `json:"save_payment_method"`
	PaymentMethodData paymentMethod       `json:"payment_method_data"`
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
			Value:    convert.FloatWithoutTrailingZeroes(request.Amount.Value, amountPrecision),
			Currency: request.Amount.Currency,
		},
		Capture:           request.Capture,
		SavePaymentMethod: request.SavePaymentMethod,
		PaymentMethodData: paymentMethod{
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

	reqURL, err := url.JoinPath(c.baseURL, paymentsEndpoint)
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
		var resp errorResponse
		if err = json.Unmarshal(respPayload, &resp); err != nil {
			err = errorFromUnresolvedResponse(reqMethod, reqURL, httpResponse.StatusCode, respPayload)
			return response, err
		}

		err = errorResponseData(resp)
		return response, err
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

type paymentCancellation struct {
	Party  string `json:"party"`
	Reason string `json:"reason"`
}

type getPaymentResponse struct {
	ID            string              `json:"id"`
	Status        string              `json:"status"`
	CapturedAt    string              `json:"captured_at"`
	IncomeAmount  paymentAmount       `json:"income_amount"`
	Cancellation  paymentCancellation `json:"cancellation_details"`
	PaymentMethod paymentMethod       `json:"payment_method"`
}

func (c *Client) GetPayment(ctx context.Context, paymentID string) (data.GetPaymentResponse, error) {
	var response data.GetPaymentResponse

	reqURL, err := url.JoinPath(c.baseURL, paymentsEndpoint, paymentID)
	if err != nil {
		return response, fmt.Errorf("failed to join base URL with get payment request endpoint: %w", err)
	}

	reqMethod := http.MethodGet

	httpRequest, err := http.NewRequestWithContext(ctx, reqMethod, reqURL, nil)
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
		var resp errorResponse
		if err = json.Unmarshal(respPayload, &resp); err != nil {
			err = errorFromUnresolvedResponse(reqMethod, reqURL, httpResponse.StatusCode, respPayload)
			return response, err
		}

		if httpResponse.StatusCode == http.StatusNotFound && resp.Code == errorResponseCodeNotFound {
			return response, data.ErrPaymentNotFound
		}

		err = errorResponseData(resp)
		return response, err
	}

	var resp getPaymentResponse
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

	if resp.IncomeAmount.Value != "" {
		response.IncomeAmount.Value, err = strconv.ParseFloat(resp.IncomeAmount.Value, 64)
		if err != nil {
			return response, fmt.Errorf("parse income_amount.value %q as float: %w", resp.IncomeAmount.Value, err)
		}
	}
	response.IncomeAmount.Currency = resp.IncomeAmount.Currency

	response.ID = resp.ID
	response.Status = resp.Status
	response.Cancellation = data.PaymentCancellation(resp.Cancellation)
	response.PaymentMethod = data.PaymentMethod{
		Type:  resp.PaymentMethod.Type,
		ID:    resp.PaymentMethod.ID,
		Saved: resp.PaymentMethod.Saved,
	}
	if resp.PaymentMethod.Card != nil {
		response.PaymentMethod.Card = data.PaymentMethodCard{
			First6:      resp.PaymentMethod.Card.First6,
			Last4:       resp.PaymentMethod.Card.Last4,
			ExpiryYear:  resp.PaymentMethod.Card.ExpiryYear,
			ExpiryMonth: resp.PaymentMethod.Card.ExpiryMonth,
			CardType:    resp.PaymentMethod.Card.CardType,
			IssuerName:  resp.PaymentMethod.Card.IssuerName,
		}
	}

	return response, nil
}
