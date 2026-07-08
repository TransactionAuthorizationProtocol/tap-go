package tap

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	didcomm "github.com/notabene-id/go-didcomm"
)

// UpdatePoliciesBody represents the body of a TAP UpdatePolicies message (TAIP-7).
type UpdatePoliciesBody struct {
	Context  string   `json:"@context"`
	Type     string   `json:"@type"`
	Policies []Policy `json:"policies"`
}

func (b *UpdatePoliciesBody) TAPType() string { return TypeUpdatePolicies }

// NewUpdatePoliciesMessage creates a new DIDComm message with an UpdatePolicies body.
func NewUpdatePoliciesMessage(from string, to []string, thid string, body *UpdatePoliciesBody) (*didcomm.Message, error) {
	if len(body.Policies) == 0 {
		return nil, fmt.Errorf("%w: missing policies", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeUpdatePolicies

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeUpdatePolicies,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
