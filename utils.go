/*
Copyright 2025 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package opa

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/nuclio/errors"
)

func sendHTTPRequest(ctx context.Context,
	httpClient *http.Client,
	method string,
	requestURL string,
	body []byte,
	headers map[string]string,
	cookies []*http.Cookie,
	expectedStatusCode int) ([]byte, *http.Response, error) {

	// create request object
	req, err := http.NewRequestWithContext(ctx, method, requestURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to create http request")
	}

	// attach cookies
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// attach headers
	for headerKey, headerValue := range headers {
		req.Header.Set(headerKey, headerValue)
	}

	// perform the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to send HTTP request")
	}

	// read response body
	var responseBody []byte
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close() // nolint: errcheck

		responseBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed to read response body")
		}
	}

	// validate status code is as expected
	if expectedStatusCode != 0 && resp != nil && resp.StatusCode != expectedStatusCode {
		return responseBody, resp, errors.Errorf(
			"Got unexpected response status code: %d. Expected: %d",
			resp.StatusCode,
			expectedStatusCode)
	}

	return responseBody, resp, nil
}

// retryUntilSuccessful retries a callback function until it returns true or timeout is reached.
// It waits for the specified interval between retries.
// Returns an error if the timeout duration is exceeded without success.
func retryUntilSuccessful(duration time.Duration, interval time.Duration, callback func() bool) error {
	timeout := time.After(duration)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Try immediately first
	if callback() {
		return nil
	}

	for {
		select {
		case <-timeout:
			return errors.New("Retry timeout exceeded")
		case <-ticker.C:
			if callback() {
				return nil
			}
		}
	}
}
