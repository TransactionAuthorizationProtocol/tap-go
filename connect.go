package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// ConnectBody represents the body of a TAP Connect message (TAIP-15).
type ConnectBody struct {
	Context     string                  `json:"@context"`
	Type        string                  `json:"@type"`
	Requester   *Party                  `json:"requester"`
	Principal   *Party                  `json:"principal"`
	Agents      []Agent                 `json:"agents"`
	Constraints *TransactionConstraints `json:"constraints"`
	Agreement   string                  `json:"agreement,omitempty"`
	Expiry      string                  `json:"expiry,omitempty"`
}

func (b *ConnectBody) TAPType() string { return TypeConnect }

// NewConnectMessage creates a new DIDComm message with a Connect body.
func NewConnectMessage(from string, to []string, body *ConnectBody) (*didcomm.Message, error) {
	if body.Requester == nil {
		return nil, fmt.Errorf("%w: missing requester", ErrInvalidBody)
	}
	if body.Principal == nil {
		return nil, fmt.Errorf("%w: missing principal", ErrInvalidBody)
	}
	if len(body.Agents) == 0 {
		return nil, fmt.Errorf("%w: missing agents", ErrInvalidBody)
	}
	if body.Constraints == nil {
		return nil, fmt.Errorf("%w: missing constraints", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeConnect

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeConnect,
		From: from,
		To:   to,
		Body: rawBody,
	}, nil
}
