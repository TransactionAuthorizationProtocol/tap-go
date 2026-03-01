package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// SettleBody represents the body of a TAP Settle message (TAIP-4).
type SettleBody struct {
	Context           string `json:"@context"`
	Type              string `json:"@type"`
	SettlementAddress string `json:"settlementAddress"`
	SettlementID      string `json:"settlementId,omitempty"`
	Amount            string `json:"amount,omitempty"`
}

func (b *SettleBody) TAPType() string { return TypeSettle }

// NewSettleMessage creates a new DIDComm message with a Settle body.
func NewSettleMessage(from string, to []string, thid string, body *SettleBody) (*didcomm.Message, error) {
	if body.SettlementAddress == "" {
		return nil, fmt.Errorf("%w: missing settlementAddress", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeSettle

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeSettle,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
