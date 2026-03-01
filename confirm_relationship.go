package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// ConfirmRelationshipBody represents the body of a TAP ConfirmRelationship message (TAIP-9).
type ConfirmRelationshipBody struct {
	Context    string        `json:"@context"`
	Type       string        `json:"@type"`
	Relationship *Relationship `json:"relationship"`
	Status     string        `json:"status"`
	ValidFrom  string        `json:"validFrom,omitempty"`
	ValidUntil string        `json:"validUntil,omitempty"`
	Details    any           `json:"details,omitempty"`
}

func (b *ConfirmRelationshipBody) TAPType() string { return TypeConfirmRelationship }

// NewConfirmRelationshipMessage creates a new DIDComm message with a ConfirmRelationship body.
func NewConfirmRelationshipMessage(from string, to []string, thid string, body *ConfirmRelationshipBody) (*didcomm.Message, error) {
	if body.Relationship == nil {
		return nil, fmt.Errorf("%w: missing relationship", ErrInvalidBody)
	}
	if body.Status == "" {
		return nil, fmt.Errorf("%w: missing status", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeConfirmRelationship

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeConfirmRelationship,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
