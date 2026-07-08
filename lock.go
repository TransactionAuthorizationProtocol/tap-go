package tap

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	didcomm "github.com/notabene-id/go-didcomm"
)

// LockBody represents the body of a TAP Lock message (TAIP-17).
type LockBody struct {
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

func (b *LockBody) TAPType() string { return TypeLock }

// NewLockMessage creates a new DIDComm message with a Lock body.
func NewLockMessage(from string, to []string, body *LockBody) (*didcomm.Message, error) {
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
	body.Type = TypeLock

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeLock,
		From: from,
		To:   to,
		Body: rawBody,
	}, nil
}
