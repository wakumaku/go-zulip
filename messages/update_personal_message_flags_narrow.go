package messages

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/narrow"
)

type UpdatePersonalMessageFlagsNarrow struct {
	zulip.APIResponseBase
	updatePersonalMessageFlagsNarrowData
}

type updatePersonalMessageFlagsNarrowData struct {
	ProcessedCount   int  `json:"processed_count"`
	UpdatedCount     int  `json:"updated_count"`
	FirstProcessedID int  `json:"first_processed_id"`
	LastProcessedID  int  `json:"last_processed_id"`
	FoundOldest      bool `json:"found_oldest"`
	FoundNewest      bool `json:"found_newest"`
}

func (g *UpdatePersonalMessageFlagsNarrow) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &g.APIResponseBase); err != nil {
		return err
	}

	if err := json.Unmarshal(b, &g.updatePersonalMessageFlagsNarrowData); err != nil {
		return err
	}

	return nil
}

type updatePersonalMessageFlagsNarrowOptions struct {
	IncludeAnchor struct {
		fieldName string
		value     bool
	}
}

type UpdatePersonalMessageFlagsNarrowOption func(*updatePersonalMessageFlagsNarrowOptions)

func UpdatePersonalMessageFlagsNarrowIncludeAnchor() UpdatePersonalMessageFlagsNarrowOption {
	return func(o *updatePersonalMessageFlagsNarrowOptions) {
		o.IncludeAnchor.fieldName = "include_anchor"
		o.IncludeAnchor.value = true
	}
}

func (svc *Service) UpdatePersonalMessageFlagsNarrow(
	ctx context.Context,
	anchor string,
	numBefore int,
	numAfter int,
	narrow narrow.Filter,
	op Operation,
	flag Flag,
	opts ...UpdatePersonalMessageFlagsNarrowOption,
) (*UpdatePersonalMessageFlagsNarrow, error) {
	const (
		method = http.MethodPost
		path   = "/api/v1/messages/flags/narrow"
	)

	narrowJSON, err := json.Marshal(narrow)
	if err != nil {
		return nil, err
	}

	msg := map[string]any{
		"anchor":     anchor,
		"num_before": numBefore,
		"num_after":  numAfter,
		"narrow":     string(narrowJSON),
		"op":         op,
		"flag":       flag,
	}

	options := updatePersonalMessageFlagsNarrowOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.IncludeAnchor.value {
		msg[options.IncludeAnchor.fieldName] = options.IncludeAnchor.value
	}

	resp := UpdatePersonalMessageFlagsNarrow{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
