package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	headerIdempotenceKey = "Idempotence-Key"
	headerContentType    = "Content-Type"
)

func (c *Client) setPaymentRequiredHeaders(req *http.Request) {
	req.SetBasicAuth(c.shopID, c.paymentsSecretKey)

	if req.Method == http.MethodPost {
		req.Header.Set(headerIdempotenceKey, generateXRequestID())
		req.Header.Set(headerContentType, "application/json")
	}
}

func (c *Client) setPayoutRequiredHeaders(req *http.Request) {
	req.SetBasicAuth(c.agentID, c.payoutsSecretKey)

	if req.Method == http.MethodPost {
		req.Header.Set(headerIdempotenceKey, generateXRequestID())
		req.Header.Set(headerContentType, "application/json")
	}
}

func generateXRequestID() string {
	uuidWithHyphen := uuid.NewString()
	return strings.Replace(uuidWithHyphen, "-", "", -1)
}
