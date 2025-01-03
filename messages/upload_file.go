package messages

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

type UploadFileResponse struct {
	zulip.APIResponseBase
	uploadFileResponseData
}

type uploadFileResponseData struct {
	FileName string `json:"filename"`
	URI      string `json:"uri"`
	URL      string `json:"url"`
}

func (aer *UploadFileResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &aer.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &aer.uploadFileResponseData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) UploadFile(ctx context.Context, filePath string) (*UploadFileResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/user_uploads"
	)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	resp := UploadFileResponse{}
	if err := svc.client.DoFileRequest(ctx, method, path, file.Name(), file, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (svc *Service) UploadFileFromBytes(ctx context.Context, fileName string, fileBytes []byte) (*UploadFileResponse, error) {
	return svc.UploadFileFromReader(ctx, fileName, bytes.NewReader(fileBytes))
}

func (svc *Service) UploadFileFromReader(ctx context.Context, fileName string, fileReader io.Reader) (*UploadFileResponse, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/user_uploads"
	)

	resp := UploadFileResponse{}
	if err := svc.client.DoFileRequest(ctx, method, path, fileName, fileReader, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
