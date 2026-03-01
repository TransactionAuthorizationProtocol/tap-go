package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// RejectBody represents the body of a TAP Reject message (TAIP-4).
type RejectBody struct {
	Context string `json:"@context"`
	Type    string `json:"@type"`
	Reason  string `json:"reason,omitempty"`
}

func (b *RejectBody) TAPType() string { return TypeReject }

// NewRejectMessage creates a new DIDComm message with a Reject body.
func NewRejectMessage(from string, to []string, thid string, body *RejectBody) (*didcomm.Message, error) {
	body.Context = TAPContext
	body.Type = TypeReject

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeReject,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
