// Package intranet provides a client for accessing the MiniEye Intranet RESTful API.
package intranet

import (
	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/services"
)

// Client represents a client for the MiniEye Intranet API.
type Client struct {
	client    *client.HTTPClient
	User      services.UserService
	Connector services.ConnectorService
	ApiKey    services.ApiKeyService
}

// NewClient creates a new MiniEye Intranet API client with the given options.
func NewClient(options ...Option) (*Client, error) {
	cfg := defaultConfig()
	for _, opt := range options {
		opt(cfg)
	}

	httpClient, err := client.NewHTTPClient(cfg)
	if err != nil {
		return nil, err
	}

	// Create user service
	userService := services.NewUserService(httpClient)

	// Create connector service
	connectorService := services.NewConnectorService(httpClient)

	// Create apikey service
	apiKeyService := services.NewApiKeyService(httpClient)

	return &Client{
		client:    httpClient,
		User:      userService,
		Connector: connectorService,
		ApiKey:    apiKeyService,
	}, nil
}

// Option is a function for configuring the client.
type Option func(*client.Config)

// Config is an alias to client.Config for backward compatibility
type Config client.Config

// defaultConfig returns a default configuration.
func defaultConfig() *client.Config {
	return &client.Config{
		BaseURL:   "https://intranet.minieye.tech/sys-api",
		UserAgent: "minieye-intranet-sdk/1.0",
	}
}

// WithBaseURL sets the base URL for API requests.
func WithBaseURL(url string) Option {
	return func(c *client.Config) {
		c.BaseURL = url
	}
}

// WithAPIKey sets the API key for authentication.
func WithAPIKey(apiKey string) Option {
	return func(c *client.Config) {
		c.APIKey = apiKey
	}
}

// WithUserAgent sets the user agent for API requests.
func WithUserAgent(userAgent string) Option {
	return func(c *client.Config) {
		c.UserAgent = userAgent
	}
}

// WithAccessKeyID sets the access key ID for STS authentication.
func WithAccessKeyID(accessKeyID string) Option {
	return func(c *client.Config) {
		c.AccessKeyID = accessKeyID
	}
}

// WithAccessKeySecret sets the access key secret for STS authentication.
func WithAccessKeySecret(accessKeySecret string) Option {
	return func(c *client.Config) {
		c.AccessKeySecret = accessKeySecret
	}
}
