package tap

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	didcomm "github.com/notabene-id/go-didcomm"
)

// RevertBody represents the body of a TAP Revert message (TAIP-4).
type RevertBody struct {
	Context           string `json:"@context"`
	Type              string `json:"@type"`
	SettlementAddress string `json:"settlementAddress"`
	Reason            string `json:"reason"`
}

func (b *RevertBody) TAPType() string { return TypeRevert }

// NewRevertMessage creates a new DIDComm message with a Revert body.
func NewRevertMessage(from string, to []string, thid string, body *RevertBody) (*didcomm.Message, error) {
	if body.SettlementAddress == "" {
		return nil, fmt.Errorf("%w: missing settlementAddress", ErrInvalidBody)
	}
	if body.Reason == "" {
		return nil, fmt.Errorf("%w: missing reason", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeRevert

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeRevert,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
