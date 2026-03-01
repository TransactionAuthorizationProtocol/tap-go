package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// TransferBody represents the body of a TAP Transfer message (TAIP-3).
type TransferBody struct {
	Context      string  `json:"@context"`
	Type         string  `json:"@type"`
	Asset        string  `json:"asset"`
	Amount       string  `json:"amount,omitempty"`
	Originator   *Party  `json:"originator,omitempty"`
	Beneficiary  *Party  `json:"beneficiary,omitempty"`
	Agents       []Agent `json:"agents"`
	SettlementID string  `json:"settlementId,omitempty"`
	Memo         string  `json:"memo,omitempty"`
	Expiry       string  `json:"expiry,omitempty"`
}

func (b *TransferBody) TAPType() string { return TypeTransfer }

// NewTransferMessage creates a new DIDComm message with a Transfer body.
func NewTransferMessage(from string, to []string, body *TransferBody) (*didcomm.Message, error) {
	if body.Asset == "" {
		return nil, fmt.Errorf("%w: missing asset", ErrInvalidBody)
	}
	if len(body.Agents) == 0 {
		return nil, fmt.Errorf("%w: missing agents", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeTransfer

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeTransfer,
		From: from,
		To:   to,
		Body: rawBody,
	}, nil
}
