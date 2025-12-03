// Package client provides the HTTP client for the MiniEye Intranet API.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/hujia-team/intranet-sdk/utils"
)

// HTTPClient represents the HTTP client for the API.
type HTTPClient struct {
	client    *http.Client
	config    *Config
	authToken string
	mu        sync.RWMutex
}

// NewHTTPClient creates a new HTTP client.
func NewHTTPClient(config *Config) (*HTTPClient, error) {
	if config.BaseURL == "" {
		return nil, utils.NewInvalidInputError("base URL is required", nil)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &HTTPClient{
		client: client,
		config: config,
	}, nil
}

// SetAuthToken sets the authentication token for the client.
func (c *HTTPClient) SetAuthToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.authToken = token
}

// GetAuthToken returns the current authentication token.
func (c *HTTPClient) GetAuthToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.authToken
}

// Config holds the configuration for the HTTP client.
type Config struct {
	BaseURL         string
	APIKey          string
	UserAgent       string
	HTTPClient      interface{}
	AccessKeyID     string
	AccessKeySecret string
}

// Do sends an HTTP request and returns the response.
func (c *HTTPClient) Do(method, endpoint string, body interface{}, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.config.BaseURL, endpoint)

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			utils.Debug("Failed to marshal request body: %v", err)
			return utils.NewInternalError("failed to marshal request body", err)
		}
		utils.Trace("Request body: %s", string(jsonBytes))
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		utils.Debug("Failed to create request: %v", err)
		return utils.NewInternalError("failed to create request", err)
	}

	utils.Trace("%s %s", method, url)

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.config.UserAgent != "" {
		req.Header.Set("User-Agent", c.config.UserAgent)
	}
	// 优先级: authToken (从登录) > STS认证 > apiKey
	authToken := c.GetAuthToken()

	if authToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	} else if c.config.AccessKeyID != "" && c.config.AccessKeySecret != "" {
		// 使用STS认证方式
		req.Header.Set("x-sts-uid", c.config.AccessKeyID)
		req.Header.Set("x-sts-token", utils.GenerateToken(c.config.AccessKeyID, c.config.AccessKeySecret))
	} else if c.config.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		utils.Debug("Failed to send request: %v", err)
		return utils.NewNetworkError("failed to send request", err)
	}
	defer resp.Body.Close()

	utils.Trace("Response status: %d", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Debug("Failed to read response body: %v", err)
		return utils.NewInternalError("failed to read response body", err)
	}

	utils.Trace("Response body: %s", string(respBody))

	// Check if the response status code is successful
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		utils.Error("API error: status=%d, body=%s", resp.StatusCode, string(respBody))

		// Map HTTP status codes to SDK error types
		var code utils.ErrorCode
		switch {
		case resp.StatusCode == 401:
			code = utils.ErrCodeUnauthorized
		case resp.StatusCode == 403:
			code = utils.ErrCodeForbidden
		case resp.StatusCode == 404:
			code = utils.ErrCodeNotFound
		case resp.StatusCode >= 500:
			code = utils.ErrCodeAPIError
		default:
			code = utils.ErrCodeAPIError
		}

		errorMsg := fmt.Sprintf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
		return utils.NewSDKError(code, errorMsg, nil)
	}

	// Parse the response body if result is provided
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			utils.Debug("Failed to unmarshal response body: %v", err)
			return utils.NewInternalError("failed to unmarshal response body", err)
		}
	}

	return nil
}

// Get sends a GET request to the specified endpoint.
func (c *HTTPClient) Get(endpoint string, result interface{}) error {
	return c.Do("GET", endpoint, nil, result)
}

// Post sends a POST request to the specified endpoint with the given body.
func (c *HTTPClient) Post(endpoint string, body interface{}, result interface{}) error {
	return c.Do("POST", endpoint, body, result)
}

// Put sends a PUT request to the specified endpoint with the given body.
func (c *HTTPClient) Put(endpoint string, body interface{}, result interface{}) error {
	return c.Do("PUT", endpoint, body, result)
}

// Delete sends a DELETE request to the specified endpoint.
func (c *HTTPClient) Delete(endpoint string, result interface{}) error {
	return c.Do("DELETE", endpoint, nil, result)
}
