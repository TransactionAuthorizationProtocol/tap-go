package tap

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	didcomm "github.com/notabene-id/go-didcomm"
)

// QuoteBody represents the body of a TAP Quote message (TAIP-18).
type QuoteBody struct {
	Context    string  `json:"@context"`
	Type       string  `json:"@type"`
	FromAsset  string  `json:"fromAsset"`
	ToAsset    string  `json:"toAsset"`
	FromAmount string  `json:"fromAmount"`
	ToAmount   string  `json:"toAmount"`
	Provider   *Party  `json:"provider"`
	Agents     []Agent `json:"agents"`
	ExpiresAt  string  `json:"expiresAt"`
}

func (b *QuoteBody) TAPType() string { return TypeQuote }

// NewQuoteMessage creates a new DIDComm message with a Quote body.
func NewQuoteMessage(from string, to []string, thid string, body *QuoteBody) (*didcomm.Message, error) {
	if body.FromAsset == "" {
		return nil, fmt.Errorf("%w: missing fromAsset", ErrInvalidBody)
	}
	if body.ToAsset == "" {
		return nil, fmt.Errorf("%w: missing toAsset", ErrInvalidBody)
	}
	if body.FromAmount == "" {
		return nil, fmt.Errorf("%w: missing fromAmount", ErrInvalidBody)
	}
	if body.ToAmount == "" {
		return nil, fmt.Errorf("%w: missing toAmount", ErrInvalidBody)
	}
	if body.Provider == nil {
		return nil, fmt.Errorf("%w: missing provider", ErrInvalidBody)
	}
	if len(body.Agents) == 0 {
		return nil, fmt.Errorf("%w: missing agents", ErrInvalidBody)
	}
	if body.ExpiresAt == "" {
		return nil, fmt.Errorf("%w: missing expiresAt", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeQuote

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeQuote,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
