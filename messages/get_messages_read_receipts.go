package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wakumaku/go-zulip"
)

type GetMessagesReadReceipts struct {
	zulip.APIResponseBase
	getMessagesReadReceiptsData
}

type getMessagesReadReceiptsData struct {
	UserIDs []int `json:"user_ids"`
}

func (g *GetMessagesReadReceipts) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.getMessagesReadReceiptsData); err != nil {
		return err
	}

	return nil
}

func (svc *Service) GetMessagesReadReceipts(ctx context.Context, messageID int) (*GetMessagesReadReceipts, error) {
	const (
		method = http.MethodGet
		path   = "/api/v1/messages/{message_id}/read_receipts"
	)
	patchPath := strings.Replace(path, "{message_id}", fmt.Sprintf("%d", messageID), 1)

	resp := GetMessagesReadReceipts{}
	if err := svc.client.DoRequest(ctx, method, patchPath, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
