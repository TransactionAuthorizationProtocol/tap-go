package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// CancelBody represents the body of a TAP Cancel message (TAIP-4).
type CancelBody struct {
	Context string `json:"@context"`
	Type    string `json:"@type"`
	By      string `json:"by"`
	Reason  string `json:"reason,omitempty"`
}

func (b *CancelBody) TAPType() string { return TypeCancel }

// NewCancelMessage creates a new DIDComm message with a Cancel body.
func NewCancelMessage(from string, to []string, thid string, body *CancelBody) (*didcomm.Message, error) {
	if body.By == "" {
		return nil, fmt.Errorf("%w: missing by", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeCancel

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeCancel,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
