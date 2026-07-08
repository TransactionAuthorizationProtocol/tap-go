package tap

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	didcomm "github.com/notabene-id/go-didcomm"
)

// CaptureBody represents the body of a TAP Capture message (TAIP-17).
type CaptureBody struct {
	Context           string `json:"@context"`
	Type              string `json:"@type"`
	Amount            string `json:"amount,omitempty"`
	SettlementAddress string `json:"settlementAddress,omitempty"`
}

func (b *CaptureBody) TAPType() string { return TypeCapture }

// NewCaptureMessage creates a new DIDComm message with a Capture body.
func NewCaptureMessage(from string, to []string, thid string, body *CaptureBody) (*didcomm.Message, error) {
	body.Context = TAPContext
	body.Type = TypeCapture

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeCapture,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
