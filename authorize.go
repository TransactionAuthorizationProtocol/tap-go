package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// AuthorizeBody represents the body of a TAP Authorize message (TAIP-4).
type AuthorizeBody struct {
	Context           string `json:"@context"`
	Type              string `json:"@type"`
	SettlementAddress string `json:"settlementAddress,omitempty"`
	SettlementAsset   string `json:"settlementAsset,omitempty"`
	Amount            string `json:"amount,omitempty"`
	Expiry            string `json:"expiry,omitempty"`
	// ApprovedTypes echoes the approved connection types on a TAIP-15
	// connection Authorize (a subset of the Connect's connectionTypes;
	// TAIPs#53). Absent on transaction authorizations.
	ApprovedTypes []string `json:"approvedTypes,omitempty"`
}

func (b *AuthorizeBody) TAPType() string { return TypeAuthorize }

// NewAuthorizeMessage creates a new DIDComm message with an Authorize body.
func NewAuthorizeMessage(from string, to []string, thid string, body *AuthorizeBody) (*didcomm.Message, error) {
	body.Context = TAPContext
	body.Type = TypeAuthorize

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeAuthorize,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
