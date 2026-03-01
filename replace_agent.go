package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// ReplaceAgentBody represents the body of a TAP ReplaceAgent message (TAIP-5).
type ReplaceAgentBody struct {
	Context     string `json:"@context"`
	Type        string `json:"@type"`
	Original    string `json:"original"`
	Replacement *Agent `json:"replacement"`
}

func (b *ReplaceAgentBody) TAPType() string { return TypeReplaceAgent }

// NewReplaceAgentMessage creates a new DIDComm message with a ReplaceAgent body.
func NewReplaceAgentMessage(from string, to []string, thid string, body *ReplaceAgentBody) (*didcomm.Message, error) {
	if body.Original == "" {
		return nil, fmt.Errorf("%w: missing original", ErrInvalidBody)
	}
	if body.Replacement == nil {
		return nil, fmt.Errorf("%w: missing replacement", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeReplaceAgent

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeReplaceAgent,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
