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

func (c *Client) setRequiredHeaders(req *http.Request) {
	req.SetBasicAuth(c.shopID, c.secretKey)
	req.Header.Set(headerIdempotenceKey, generateXRequestID())

	if req.Method == http.MethodPost {
		req.Header.Set(headerContentType, "application/json")
	}
}

func generateXRequestID() string {
	uuidWithHyphen := uuid.NewString()
	return strings.Replace(uuidWithHyphen, "-", "", -1)
}
