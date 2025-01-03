package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/realtime/events"
)

type GetEventsEventQueueResponse struct {
	zulip.APIResponseBase
	getEventsEventQueueData
}

type getEventsEventQueueData struct {
	Events []events.Event
}

func (g *GetEventsEventQueueResponse) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &g.APIResponseBase); err != nil {
		return err
	}

	eventsMap := struct {
		Events []map[string]any `json:"events"`
	}{}

	if err := json.Unmarshal(data, &eventsMap); err != nil {
		return err
	}

	for _, e := range eventsMap.Events {

		var itemType events.EventType
		if etype, ok := e["type"]; ok {
			it, ok := etype.(string)
			if !ok {
				return errors.New("type is not a string")
			}
			itemType = events.EventType(it)
		} else {
			return errors.New("type field not found")
		}

		var ev events.Event
		switch itemType {
		case events.AlertWordsType:
			ev = &events.AlertWords{}
		case events.HeartbeatType:
			ev = &events.Heartbeat{}
		case events.MessageType:
			ev = &events.Message{}
		case events.AttachmentType:
			ev = &events.Attachment{}
		case events.PresenceType:
			ev = &events.Presence{}
		case events.RealmEmojiType:
			ev = &events.RealmEmoji{}
		case events.RealmUserType:
			ev = &events.RealmUser{}
		case events.SubmessageType:
			ev = &events.Submessage{}
		case events.TypingType:
			ev = &events.Typing{}
		case events.UpdateMessageType:
			ev = &events.UpdateMessage{}
		default:
			ev = &events.Unknown{}
		}

		itemData, err := json.Marshal(e)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(itemData, ev); err != nil {
			return err
		}

		g.Events = append(g.Events, ev)
	}

	return nil
}

type getEventsEventQueueOptions struct {
	lastEventID int
	dontBlock   bool
}

type GetEventsEventQueueOption func(*getEventsEventQueueOptions)

func LastEventID(id int) GetEventsEventQueueOption {
	return func(geeqo *getEventsEventQueueOptions) {
		geeqo.lastEventID = id
	}
}

func DontBlock() GetEventsEventQueueOption {
	return func(geeqo *getEventsEventQueueOptions) {
		geeqo.dontBlock = true
	}
}

func (svc *Service) GetEventsEventQueue(ctx context.Context, eventQueueID string, options ...GetEventsEventQueueOption) (*GetEventsEventQueueResponse, error) {
	const (
		method = http.MethodGet
		path   = "/api/v1/events"
	)

	opts := getEventsEventQueueOptions{
		lastEventID: -1,
		dontBlock:   false,
	}
	for _, opt := range options {
		opt(&opts)
	}

	msg := map[string]any{
		"queue_id":      eventQueueID,
		"last_event_id": opts.lastEventID,
		"dont_block":    opts.dontBlock,
	}

	resp := GetEventsEventQueueResponse{}
	if err := svc.client.DoRequest(ctx, method, path, msg, &resp, zulip.WithTimeout(zulip.RESTClientLongPollTimeout)); err != nil {
		return nil, err
	}

	return &resp, nil
}
