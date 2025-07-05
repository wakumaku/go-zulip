package org_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip/org"
)

func TestUploadCustomEmoji(t *testing.T) {
	client := createMockClient(`{
		"msg": "",
		"result": "success"
	}`)

	service := org.NewService(client)

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test-emoji-*.png")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write some test data to the file
	_, err = tempFile.Write([]byte("test emoji data"))
	assert.NoError(t, err)
	tempFile.Close()

	resp, err := service.UploadCustomEmoji(context.Background(), "test-emoji", tempFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Result())
}
