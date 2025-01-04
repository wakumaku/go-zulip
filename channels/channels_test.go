package channels_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/wakumaku/go-zulip"
)

// mockClient is a mock implementation of zulip.RESTClient
// just for testing purposes, cannot be used concurrently on the same instance
type mockClient struct {
	response   string
	method     string
	path       string
	paramsSent map[string]any // sort of spy for testing input parameters
}

func (mc *mockClient) DoRequest(ctx context.Context, method, path string, data map[string]any, response zulip.APIResponse, opts ...zulip.DoRequestOption) error {
	mc.method = method
	mc.path = path
	mc.paramsSent = data
	return json.Unmarshal([]byte(mc.response), response)
}

func (mc *mockClient) DoFileRequest(ctx context.Context, method, path string, fileName string, file io.Reader, response zulip.APIResponse, opts ...zulip.DoRequestOption) error {
	return errors.New("not implemented")
}

// createMockClient creates a mockClient with the given response
// TODO: other complex behaviours: 4xx, timeouts, etc.
func createMockClient(response string) zulip.RESTClient {
	return &mockClient{
		response: response,
	}
}
