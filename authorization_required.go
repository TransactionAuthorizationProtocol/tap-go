package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// AuthorizationRequiredBody represents the body of a TAP AuthorizationRequired message (TAIP-4).
type AuthorizationRequiredBody struct {
	Context          string `json:"@context"`
	Type             string `json:"@type"`
	AuthorizationURL string `json:"authorizationUrl"`
	Expires          string `json:"expires"`
	From             string `json:"from,omitempty"`
}

func (b *AuthorizationRequiredBody) TAPType() string { return TypeAuthorizationRequired }

// NewAuthorizationRequiredMessage creates a new DIDComm message with an AuthorizationRequired body.
func NewAuthorizationRequiredMessage(from string, to []string, thid string, body *AuthorizationRequiredBody) (*didcomm.Message, error) {
	if body.AuthorizationURL == "" {
		return nil, fmt.Errorf("%w: missing authorizationUrl", ErrInvalidBody)
	}
	if body.Expires == "" {
		return nil, fmt.Errorf("%w: missing expires", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeAuthorizationRequired

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeAuthorizationRequired,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
