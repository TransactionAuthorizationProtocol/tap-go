package tap

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	didcomm "github.com/notabene-id/go-didcomm"
)

// AddAgentsBody represents the body of a TAP AddAgents message (TAIP-5).
type AddAgentsBody struct {
	Context string  `json:"@context"`
	Type    string  `json:"@type"`
	Agents  []Agent `json:"agents"`
}

func (b *AddAgentsBody) TAPType() string { return TypeAddAgents }

// NewAddAgentsMessage creates a new DIDComm message with an AddAgents body.
func NewAddAgentsMessage(from string, to []string, thid string, body *AddAgentsBody) (*didcomm.Message, error) {
	if len(body.Agents) == 0 {
		return nil, fmt.Errorf("%w: missing agents", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeAddAgents

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeAddAgents,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
