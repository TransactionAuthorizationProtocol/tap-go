package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// PaymentBody represents the body of a TAP Payment message (TAIP-14).
type PaymentBody struct {
	Context                    string   `json:"@context"`
	Type                       string   `json:"@type"`
	Amount                     string   `json:"amount"`
	Asset                      string   `json:"asset,omitempty"`
	Currency                   string   `json:"currency,omitempty"`
	Merchant                   *Party   `json:"merchant"`
	Customer                   *Party   `json:"customer,omitempty"`
	Agents                     []Agent  `json:"agents"`
	SupportedAssets            []any    `json:"supportedAssets,omitempty"`
	FallbackSettlementAddresses []string `json:"fallbackSettlementAddresses,omitempty"`
	Expiry                     string   `json:"expiry,omitempty"`
	Invoice                    any      `json:"invoice,omitempty"`
	Policies                   []Policy `json:"policies,omitempty"`
}

func (b *PaymentBody) TAPType() string { return TypePayment }

// NewPaymentMessage creates a new DIDComm message with a Payment body.
func NewPaymentMessage(from string, to []string, body *PaymentBody) (*didcomm.Message, error) {
	if body.Amount == "" {
		return nil, fmt.Errorf("%w: missing amount", ErrInvalidBody)
	}
	if body.Merchant == nil {
		return nil, fmt.Errorf("%w: missing merchant", ErrInvalidBody)
	}
	if len(body.Agents) == 0 {
		return nil, fmt.Errorf("%w: missing agents", ErrInvalidBody)
	}
	if body.Asset == "" && body.Currency == "" {
		return nil, fmt.Errorf("%w: missing asset or currency", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypePayment

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypePayment,
		From: from,
		To:   to,
		Body: rawBody,
	}, nil
}
