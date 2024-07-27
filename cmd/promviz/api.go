package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

type PrometheusClient struct {
	api v1.API
}

func NewPrometheusClient(baseURL string) (*PrometheusClient, error) {
	httpClient := &http.Client{
		Transport: LoggingRoundTripper{Proxied: http.DefaultTransport},
	}

	client, err := api.NewClient(api.Config{
		Address: baseURL,
		Client:  httpClient,
	})
	if err != nil {
		return nil, err
	}

	return &PrometheusClient{
		api: v1.NewAPI(client),
	}, nil
}

func (c *PrometheusClient) Query(ctx context.Context, query string, ts time.Time) (model.Value, error) {
	result, warnings, err := c.api.Query(ctx, query, ts)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	if len(warnings) > 0 {
		// Log warnings if needed
	}

	return result, nil
}

type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the request
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n\nREQUEST:\n%s\n", string(reqDump))

	// Perform the request
	resp, err := lrt.Proxied.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Log the response
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	fmt.Printf("RESPONSE:\n%s\n\n\n", string(respDump))

	// Create a new response with a copy of the body
	bodyCopy := &bytes.Buffer{}
	io.Copy(bodyCopy, resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bodyCopy)

	return resp, nil
}
