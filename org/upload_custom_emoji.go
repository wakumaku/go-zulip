package org

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/wakumaku/go-zulip"
)

type UploadCustomEmojiResponse struct {
	zulip.APIResponseBase
}

func (u *UploadCustomEmojiResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &u.APIResponseBase); err != nil {
		return err
	}

	return nil
}

func (svc *Service) UploadCustomEmoji(ctx context.Context, name, filePath string) (*UploadCustomEmojiResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/realm/emoji/%s"
	)
	patchPath := fmt.Sprintf(path, name)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	resp := UploadCustomEmojiResponse{}
	if err := svc.client.DoFileRequest(ctx, method, patchPath, file.Name(), file, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (svc *Service) UploadCustomEmojiFromBytes(ctx context.Context, name, fileName string, fileBytes []byte) (*UploadCustomEmojiResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/realm/emoji/%s"
	)
	patchPath := fmt.Sprintf(path, name)

	resp := UploadCustomEmojiResponse{}
	if err := svc.client.DoFileRequest(ctx, method, patchPath, fileName, bytes.NewReader(fileBytes), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (svc *Service) UploadCustomEmojiFromReader(ctx context.Context, name, fileName string, fileReader io.Reader) (*UploadCustomEmojiResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/realm/emoji/%s"
	)
	patchPath := fmt.Sprintf(path, name)

	resp := UploadCustomEmojiResponse{}
	if err := svc.client.DoFileRequest(ctx, method, patchPath, fileName, fileReader, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
