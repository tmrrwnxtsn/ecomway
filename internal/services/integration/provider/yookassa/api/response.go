package api

import (
	"fmt"
	"log/slog"
	"net/http"
)

func closeResponseBody(resp *http.Response) {
	if resp == nil {
		slog.Warn("failed to close HTTP response body", "error", "response is nil")
		return
	}
	if err := resp.Body.Close(); err != nil {
		slog.Warn("failed to close HTTP response body", "error", err)
	}
}

func errorFromUnresolvedResponse(requestMethod, requestURL string, responseStatusCode int, responsePayload []byte) error {
	if len(responsePayload) > 0 {
		return fmt.Errorf(
			"%v request to %v ended with status %v and body: %s",
			requestMethod, requestURL, responseStatusCode, responsePayload,
		)
	}
	return fmt.Errorf(
		"%v request to %v ended with status %v",
		requestMethod, requestURL, responseStatusCode,
	)
}
