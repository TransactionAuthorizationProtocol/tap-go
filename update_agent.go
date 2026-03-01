package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// UpdateAgentBody represents the body of a TAP UpdateAgent message (TAIP-5).
type UpdateAgentBody struct {
	Context     string `json:"@context"`
	Type        string `json:"@type"`
	Agent       *Agent `json:"agent"`
	PreviousDID string `json:"previousDid,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Effective   string `json:"effective,omitempty"`
}

func (b *UpdateAgentBody) TAPType() string { return TypeUpdateAgent }

// NewUpdateAgentMessage creates a new DIDComm message with an UpdateAgent body.
func NewUpdateAgentMessage(from string, to []string, thid string, body *UpdateAgentBody) (*didcomm.Message, error) {
	if body.Agent == nil {
		return nil, fmt.Errorf("%w: missing agent", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeUpdateAgent

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeUpdateAgent,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
