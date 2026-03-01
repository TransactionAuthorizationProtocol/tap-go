package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// UpdatePartyBody represents the body of a TAP UpdateParty message (TAIP-6).
type UpdatePartyBody struct {
	Context       string `json:"@context"`
	Type          string `json:"@type"`
	Party         *Party `json:"party"`
	Role          string `json:"role"`
	PreviousParty *Party `json:"previousParty,omitempty"`
	Reason        string `json:"reason,omitempty"`
	Effective     string `json:"effective,omitempty"`
}

func (b *UpdatePartyBody) TAPType() string { return TypeUpdateParty }

// NewUpdatePartyMessage creates a new DIDComm message with an UpdateParty body.
func NewUpdatePartyMessage(from string, to []string, thid string, body *UpdatePartyBody) (*didcomm.Message, error) {
	if body.Party == nil {
		return nil, fmt.Errorf("%w: missing party", ErrInvalidBody)
	}
	if body.Role == "" {
		return nil, fmt.Errorf("%w: missing role", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeUpdateParty

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeUpdateParty,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
