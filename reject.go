package tap

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	didcomm "github.com/notabene-id/go-didcomm"
)

// RejectBody represents the body of a TAP Reject message (TAIP-4).
//
// Reason is free text for people. Code is a machine-readable reason drawn from
// the ISO 20022 ExternalStatusReason1Code list — the codes TAIP-19 maps a Reject
// onto (pacs.002 RJCT) — so a receiver can act on why without parsing Reason.
// Both are optional; a receiver that does not know a code ignores it.
type RejectBody struct {
	Context string `json:"@context"`
	Type    string `json:"@type"`
	Reason  string `json:"reason,omitempty"`
	Code    string `json:"code,omitempty"`
}

// Reject codes (ISO 20022 ExternalStatusReason1Code).
const (
	// RejectCodeInvalidCreditorAccountNumber (AC03, "Creditor account number
	// invalid or missing"): the beneficiary side does not hold the settlement
	// address the transfer names.
	RejectCodeInvalidCreditorAccountNumber = "AC03"
)

func (b *RejectBody) TAPType() string { return TypeReject }

// NewRejectMessage creates a new DIDComm message with a Reject body.
func NewRejectMessage(from string, to []string, thid string, body *RejectBody) (*didcomm.Message, error) {
	body.Context = TAPContext
	body.Type = TypeReject

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeReject,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
