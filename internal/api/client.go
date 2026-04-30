package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const defaultBaseURL = "https://api.tailscale.com/api/v2"

// Client is an HTTP client for the Tailscale v2 API.
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
	debug      bool
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithDebug enables debug logging of HTTP requests to stderr.
func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

// NewClient creates a new Tailscale API client.
func NewClient(apiToken string, opts ...Option) *Client {
	c := &Client{
		baseURL:    defaultBaseURL,
		apiToken:   apiToken,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// apiError represents an error response from the Tailscale API.
type apiError struct {
	Message string `json:"message"`
}

// Get performs a GET request to the given API path.
func (c *Client) Get(path string) ([]byte, error) {
	return c.Do("GET", path, nil)
}

// Post performs a POST request to the given API path.
func (c *Client) Post(path string, body io.Reader) ([]byte, error) {
	return c.Do("POST", path, body)
}

// Put performs a PUT request to the given API path.
func (c *Client) Put(path string, body io.Reader) ([]byte, error) {
	return c.Do("PUT", path, body)
}

// Patch performs a PATCH request to the given API path.
func (c *Client) Patch(path string, body io.Reader) ([]byte, error) {
	return c.Do("PATCH", path, body)
}

// Delete performs a DELETE request to the given API path.
func (c *Client) Delete(path string) ([]byte, error) {
	return c.Do("DELETE", path, nil)
}

// Do performs an HTTP request to the Tailscale API.
func (c *Client) Do(method, path string, body io.Reader) ([]byte, error) {
	url := c.baseURL + path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s %s\n", method, url)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if c.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Status: %d\n", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr apiError
		if jsonErr := json.Unmarshal(respBody, &apiErr); jsonErr == nil && apiErr.Message != "" {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, apiErr.Message)
		}
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// DoJSON performs an HTTP request with a JSON-encoded body.
func (c *Client) DoJSON(method, path string, body interface{}) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		reader = bytes.NewReader(data)
	}
	return c.Do(method, path, reader)
}
