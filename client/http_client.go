// Package client provides the HTTP client for the MiniEye Intranet API.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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

type RawResponse struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// NewHTTPClient creates a new HTTP client.
func NewHTTPClient(config *Config) (*HTTPClient, error) {
	if config.BaseURL == "" {
		return nil, utils.NewInvalidInputError("base URL is required", nil)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	if providedClient, ok := config.HTTPClient.(*http.Client); ok && providedClient != nil {
		httpClient = providedClient
	}

	return &HTTPClient{
		client: httpClient,
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	c.applyDefaultHeaders(req)
	respBody, err := c.doRequest(req)
	if err != nil {
		return err
	}
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

// PostMultipart sends a multipart/form-data POST request.
func (c *HTTPClient) PostMultipart(endpoint string, body *bytes.Buffer, contentType string, headers map[string]string, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.config.BaseURL, endpoint)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		utils.Debug("Failed to create multipart request: %v", err)
		return utils.NewInternalError("failed to create multipart request", err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	c.applyDefaultHeaders(req)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return err
	}
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			utils.Debug("Failed to unmarshal multipart response body: %v", err)
			return utils.NewInternalError("failed to unmarshal response body", err)
		}
	}
	return nil
}

func (c *HTTPClient) PostMultipartRaw(endpoint string, body *bytes.Buffer, contentType string, headers map[string]string) (*RawResponse, error) {
	req, err := c.newRequest(http.MethodPost, fmt.Sprintf("%s%s", c.config.BaseURL, endpoint), body)
	if err != nil {
		utils.Debug("Failed to create multipart request: %v", err)
		return nil, utils.NewInternalError("failed to create multipart request", err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	c.applyDefaultHeaders(req)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return c.doRawRequest(req)
}

func (c *HTTPClient) PostRaw(endpoint string, body interface{}, headers map[string]string) (*RawResponse, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, utils.NewInternalError("failed to marshal request body", err)
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}
	req, err := c.newRequest(http.MethodPost, fmt.Sprintf("%s%s", c.config.BaseURL, endpoint), bodyReader)
	if err != nil {
		return nil, utils.NewInternalError("failed to create request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	c.applyDefaultHeaders(req)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return c.doRawRequest(req)
}

func (c *HTTPClient) PostMultipartRawURL(rawURL string, body *bytes.Buffer, contentType string, headers map[string]string) (*RawResponse, error) {
	req, err := c.newRequest(http.MethodPost, rawURL, body)
	if err != nil {
		utils.Debug("Failed to create multipart request: %v", err)
		return nil, utils.NewInternalError("failed to create multipart request", err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	c.applyDefaultHeaders(req)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return c.doRawRequest(req)
}

func (c *HTTPClient) PostRawURL(rawURL string, body interface{}, headers map[string]string) (*RawResponse, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, utils.NewInternalError("failed to marshal request body", err)
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}
	req, err := c.newRequest(http.MethodPost, rawURL, bodyReader)
	if err != nil {
		return nil, utils.NewInternalError("failed to create request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	c.applyDefaultHeaders(req)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return c.doRawRequest(req)
}

// BuildMultipartBody creates a multipart body from form fields and one file field.
func BuildMultipartBody(fields map[string]string, fileField string, fileName string, fileContent []byte) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, "", err
		}
	}
	part, err := writer.CreateFormFile(fileField, fileName)
	if err != nil {
		return nil, "", err
	}
	if _, err := part.Write(fileContent); err != nil {
		return nil, "", err
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return body, writer.FormDataContentType(), nil
}

func (c *HTTPClient) applyDefaultHeaders(req *http.Request) {
	if c.config.UserAgent != "" {
		req.Header.Set("User-Agent", c.config.UserAgent)
	}
	authToken := c.GetAuthToken()
	if authToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		return
	}
	if c.config.AccessKeyID != "" && c.config.AccessKeySecret != "" {
		req.Header.Set("x-sts-uid", c.config.AccessKeyID)
		req.Header.Set("x-sts-token", utils.GenerateToken(c.config.AccessKeyID, c.config.AccessKeySecret))
		return
	}
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	}
}

func (c *HTTPClient) doRequest(req *http.Request) ([]byte, error) {
	rawResp, err := c.doRawRequest(req)
	if err != nil {
		return nil, err
	}
	return rawResp.Body, nil
}

func (c *HTTPClient) doRawRequest(req *http.Request) (*RawResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		utils.Debug("Failed to send request: %v", err)
		return nil, utils.NewNetworkError("failed to send request", err)
	}
	defer resp.Body.Close()

	utils.Trace("Response status: %d", resp.StatusCode)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Debug("Failed to read response body: %v", err)
		return nil, utils.NewInternalError("failed to read response body", err)
	}
	utils.Trace("Response body: %s", string(respBody))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		utils.Error("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
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
		return &RawResponse{StatusCode: resp.StatusCode, Body: respBody, Header: resp.Header.Clone()}, utils.NewSDKError(code, errorMsg, nil)
	}
	return &RawResponse{StatusCode: resp.StatusCode, Body: respBody, Header: resp.Header.Clone()}, nil
}

func (c *HTTPClient) newRequest(method string, rawURL string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, rawURL, body)
}
