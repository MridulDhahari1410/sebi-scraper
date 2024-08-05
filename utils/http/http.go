package http

import (
	"context"
	"net/http"

	httpclient "github.com/angel-one/go-http-client"
	"github.com/sinhashubham95/go-utils/log"
)

// Client is the set of methods for the http client.
type Client interface {
	Request(request *httpclient.Request) (*http.Response, error)
}

var client *httpclient.Client

// InitHTTPClient is used to initialise the http client.
func InitHTTPClient(configs ...*httpclient.RequestConfig) {
	client = httpclient.ConfigureHTTPClient(configs...).WithLogger(func(ctx context.Context,
		msg string) {
		log.Info(ctx).Msg(msg)
	})
}

// NewRequestConfig is used to create a new request config.
func NewRequestConfig(name string, configs map[string]any) *httpclient.RequestConfig {
	return httpclient.NewRequestConfig(name, configs)
}

// NewRequest is used to create a new request.
func NewRequest(name string) *httpclient.Request {
	return httpclient.NewRequest(name)
}

// Get is used to get the client instance.
func Get() Client {
	return client
}
