package query

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Default retry amount
const defaultRetry = 1

// Default delay between retries
const defaultDelayBetweenRetries = 10 * time.Second

// Default timeout for HTTP Request
const defaultTimeoutInSeconds = 300

// HTTPRequest default HTTP Request structure
type HTTPRequest struct {
	URL     string
	Method  string
	Headers map[string]string
	Retry   int
	Timeout time.Duration
	Body    io.Reader
}

// HTTPResponse default query engine response
type HTTPResponse struct {
	Body  []byte
	Error error
}

// MakeHTTPRequest makes HTTP request to given `url` with `headers` and in case of error
// it will retry request `retry` amount of times. And only after it (if it's still error) error will be returned.
// Automatically timeout between requests will be calculated using `random`.
// Note for `timeout` waiting this function uses `time.Sleep()` so it will block execution flow.
// Better to be used in go-routine.
func MakeHTTPRequest(r *HTTPRequest) *HTTPResponse {
	if r == nil {
		return &HTTPResponse{
			Error: fmt.Errorf("failed to make HTTP request to `nil`"),
		}
	}

	// Check for non set Retry
	if r.Retry == 0 {
		r.Retry = defaultRetry
	}

	step := 1
	var res []byte
	var err error

	for step <= r.Retry {
		res, err = doMakeHTTPRequest(r)
		if err != nil {
			time.Sleep(defaultDelayBetweenRetries)
			step++
			continue
		}
		// All ok no `err` received
		break
	}

	return &HTTPResponse{
		Body:  res,
		Error: err,
	}
}

func doMakeHTTPRequest(r *HTTPRequest) ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("failed to make HTTP request to `nil`")
	}

	// Check default method
	if r.Method == "" {
		r.Method = "GET"
	}
	// Binding default timeout
	if r.Timeout == time.Duration(0) {
		r.Timeout = defaultTimeoutInSeconds * time.Second
	}

	client := &http.Client{
		Timeout: r.Timeout,
	}
	req, err := http.NewRequest(r.Method, r.URL, r.Body)
	if err != nil {
		return nil, err
	}
	if r.Headers != nil {
		for k, v := range r.Headers {
			req.Header.Add(k, v)
		}
	}
	//fmt.Println("doMakeHTTPRequest", r.URL)
	// Perform HTTP request
	resp, err := client.Do(req)
	//fmt.Println("doMakeHTTPRequest:resp", r.Method, r.URL, r.Headers, resp.StatusCode)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("failed to make HTTP  requestto %s, got %d status code", r.URL, resp.StatusCode)
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
