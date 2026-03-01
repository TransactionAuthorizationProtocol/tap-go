package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// RemoveAgentBody represents the body of a TAP RemoveAgent message (TAIP-5).
type RemoveAgentBody struct {
	Context string `json:"@context"`
	Type    string `json:"@type"`
	Agent   string `json:"agent"`
}

func (b *RemoveAgentBody) TAPType() string { return TypeRemoveAgent }

// NewRemoveAgentMessage creates a new DIDComm message with a RemoveAgent body.
func NewRemoveAgentMessage(from string, to []string, thid string, body *RemoveAgentBody) (*didcomm.Message, error) {
	if body.Agent == "" {
		return nil, fmt.Errorf("%w: missing agent", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeRemoveAgent

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeRemoveAgent,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
