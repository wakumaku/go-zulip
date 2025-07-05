package specialty_test

import (
	"context"
	"encoding/json"
	"io"

	"github.com/wakumaku/go-zulip"
)

// mockClient is a mock implementation of zulip.RESTClient
type mockClient struct {
	response   string
	method     string
	path       string
	paramsSent map[string]any
}

func (mc *mockClient) DoRequest(ctx context.Context, method, path string, data map[string]any, response zulip.APIResponse, opts ...zulip.DoRequestOption) error {
	mc.method = method
	mc.path = path
	mc.paramsSent = data
	return json.Unmarshal([]byte(mc.response), response)
}

func (mc *mockClient) DoFileRequest(ctx context.Context, method, path string, fileName string, file io.Reader, response zulip.APIResponse, opts ...zulip.DoRequestOption) error {
	return nil
}

func createMockClient(response string) zulip.RESTClient {
	return &mockClient{
		response: response,
	}
}
