package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// RFQBody represents the body of a TAP RFQ (Request for Quote) message (TAIP-18).
type RFQBody struct {
	Context    string   `json:"@context"`
	Type       string   `json:"@type"`
	FromAssets []string `json:"fromAssets"`
	ToAssets   []string `json:"toAssets"`
	FromAmount string   `json:"fromAmount,omitempty"`
	ToAmount   string   `json:"toAmount,omitempty"`
	Requester  *Party   `json:"requester"`
	Provider   *Party   `json:"provider,omitempty"`
	Agents     []Agent  `json:"agents"`
	Policies   []Policy `json:"policies,omitempty"`
}

func (b *RFQBody) TAPType() string { return TypeRFQ }

// NewRFQMessage creates a new DIDComm message with an RFQ body.
func NewRFQMessage(from string, to []string, body *RFQBody) (*didcomm.Message, error) {
	if len(body.FromAssets) == 0 {
		return nil, fmt.Errorf("%w: missing fromAssets", ErrInvalidBody)
	}
	if len(body.ToAssets) == 0 {
		return nil, fmt.Errorf("%w: missing toAssets", ErrInvalidBody)
	}
	if body.Requester == nil {
		return nil, fmt.Errorf("%w: missing requester", ErrInvalidBody)
	}
	if len(body.Agents) == 0 {
		return nil, fmt.Errorf("%w: missing agents", ErrInvalidBody)
	}
	if body.FromAmount == "" && body.ToAmount == "" {
		return nil, fmt.Errorf("%w: missing fromAmount or toAmount", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeRFQ

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeRFQ,
		From: from,
		To:   to,
		Body: rawBody,
	}, nil
}
