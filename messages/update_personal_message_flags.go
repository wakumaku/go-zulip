package messages

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
)

type UpdatePersonalMessageFlags struct {
	zulip.APIResponseBase
	updatePersonalMessageFlagsData
}

type updatePersonalMessageFlagsData struct {
	Messages []int `json:"messages"`
}

func (g *UpdatePersonalMessageFlags) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.updatePersonalMessageFlagsData); err != nil {
		return err
	}

	return nil
}

type Operation string

const (
	OperationAdd    Operation = "add"
	OperationRemove Operation = "remove"
)

type Flag string

const (
	FlagRead                    Flag = "read"
	FlagStarred                 Flag = "starred"
	FlagCollapsed               Flag = "collapsed"
	FlagMentioned               Flag = "mentioned"
	FlagStreamWildcardMentioned Flag = "stream_wildcard_mentioned"
	FlagTopicWildcardMentioned  Flag = "topic_wildcard_mentioned"
	FlagHasAlertWord            Flag = "has_alert_word"
	FlagHistorical              Flag = "historical"
	// FlagWilcardMentioned deprecated
	// FlagWilcardMentioned Flag = "wildcard_mentioned"
)

func (svc *Service) UpdatePersonalMessageFlags(ctx context.Context, messageIDs []int, op Operation, flag Flag) (*UpdatePersonalMessageFlags, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/messages/flags"
	)

	messageIDsJSON, err := json.Marshal(messageIDs)
	if err != nil {
		return nil, err
	}

	msg := map[string]any{
		"messages": string(messageIDsJSON),
		"op":       op,
		"flag":     flag,
	}

	resp := UpdatePersonalMessageFlags{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
