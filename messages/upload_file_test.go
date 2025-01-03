package messages_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/messages"
)

func TestUploadFile(t *testing.T) {
	client := createMockClient(`{
    "filename": "zulip.txt",
    "msg": "",
    "result": "success",
    "uri": "/user_uploads/1/4e/m2A3MSqFnWRLUf9SaPzQ0Up_/zulip.txt",
    "url": "/user_uploads/1/4e/m2A3MSqFnWRLUf9SaPzQ0Up_/zulip.txt"
}`)

	messagesSvc := messages.NewService(client)

	// create a temporary file to upload
	f, err := os.CreateTemp("", "zulip.txt")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	msg := map[string]any{
		"filename": f.Name(),
	}

	resp, err := messagesSvc.UploadFile(context.Background(),
		f.Name(),
	)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.IsSuccess())

	// becuse file name is generated, we can't compare it directly
	assert.Contains(t, f.Name(), resp.FileName)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/user_uploads", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}

func TestUploadFileFromBytes(t *testing.T) {
	client := createMockClient(`{
    "filename": "zulip.txt",
    "msg": "",
    "result": "success",
    "uri": "/user_uploads/1/4e/m2A3MSqFnWRLUf9SaPzQ0Up_/zulip.txt",
    "url": "/user_uploads/1/4e/m2A3MSqFnWRLUf9SaPzQ0Up_/zulip.txt"
}`)

	messagesSvc := messages.NewService(client)

	// create a temporary file to upload
	f, err := os.CreateTemp("", "zulip.txt")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	msg := map[string]any{
		"filename": f.Name(),
	}

	resp, err := messagesSvc.UploadFileFromBytes(context.Background(),
		f.Name(),
		[]byte("hello world"),
	)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.IsSuccess())

	// becuse file name is generated, we can't compare it directly
	assert.Contains(t, f.Name(), resp.FileName)

	// validate the parameters sent are correct
	assert.Equal(t, "/api/v1/user_uploads", client.(*mockClient).path)
	assert.Equal(t, msg, client.(*mockClient).paramsSent)
}
