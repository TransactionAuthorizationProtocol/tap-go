package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// EscrowBody represents the body of a TAP Escrow message (TAIP-17).
type EscrowBody struct {
	Context     string  `json:"@context"`
	Type        string  `json:"@type"`
	Asset       string  `json:"asset,omitempty"`
	Currency    string  `json:"currency,omitempty"`
	Amount      string  `json:"amount"`
	Originator  *Party  `json:"originator"`
	Beneficiary *Party  `json:"beneficiary"`
	Expiry      string  `json:"expiry"`
	Agents      []Agent `json:"agents"`
	Agreement   string  `json:"agreement,omitempty"`
}

func (b *EscrowBody) TAPType() string { return TypeEscrow }

// NewEscrowMessage creates a new DIDComm message with an Escrow body.
func NewEscrowMessage(from string, to []string, body *EscrowBody) (*didcomm.Message, error) {
	if body.Amount == "" {
		return nil, fmt.Errorf("%w: missing amount", ErrInvalidBody)
	}
	if body.Originator == nil {
		return nil, fmt.Errorf("%w: missing originator", ErrInvalidBody)
	}
	if body.Beneficiary == nil {
		return nil, fmt.Errorf("%w: missing beneficiary", ErrInvalidBody)
	}
	if body.Expiry == "" {
		return nil, fmt.Errorf("%w: missing expiry", ErrInvalidBody)
	}
	if len(body.Agents) == 0 {
		return nil, fmt.Errorf("%w: missing agents", ErrInvalidBody)
	}
	if body.Asset == "" && body.Currency == "" {
		return nil, fmt.Errorf("%w: missing asset or currency", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeEscrow

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeEscrow,
		From: from,
		To:   to,
		Body: rawBody,
	}, nil
}
