package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// ConfirmRelationshipBody is the body of a ConfirmRelationship message (TAIP-9):
// the Agent payload, asserting that @id acts for the entity in "for".
type ConfirmRelationshipBody struct {
	Context string   `json:"@context"`
	Type    string   `json:"@type"`
	ID      string   `json:"@id"` // the agent being confirmed
	For     ForField `json:"for"` // the entity it acts on behalf of
	Role    string   `json:"role,omitempty"`
}

func (b *ConfirmRelationshipBody) TAPType() string { return TypeConfirmRelationship }

// NewConfirmRelationshipMessage creates a new DIDComm message with a ConfirmRelationship body.
func NewConfirmRelationshipMessage(from string, to []string, thid string, body *ConfirmRelationshipBody) (*didcomm.Message, error) {
	if body.ID == "" {
		return nil, fmt.Errorf("%w: missing @id", ErrInvalidBody)
	}
	if body.For.IsEmpty() {
		return nil, fmt.Errorf("%w: missing for", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeAgent

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
