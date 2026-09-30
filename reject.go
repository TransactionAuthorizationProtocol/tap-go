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
// Both are optional. NewRejectMessage sends only the codes defined below;
// ParseBody keeps any code, so a receiver ignores one it does not know rather
// than dropping the message.
type RejectBody struct {
	Context string     `json:"@context"`
	Type    string     `json:"@type"`
	Reason  string     `json:"reason,omitempty"`
	Code    RejectCode `json:"code,omitempty"`
}

// RejectCode is an ISO 20022 ExternalStatusReason1Code value on a Reject.
type RejectCode string

// Reject codes this library sends. Add a value here to allow sending it.
const (
	// RejectCodeInvalidCreditorAccountNumber (AC03, "Creditor account number
	// invalid or missing"): the beneficiary side does not hold the settlement
	// address the transfer names.
	RejectCodeInvalidCreditorAccountNumber RejectCode = "AC03"
)

// Known reports whether c is a code this library sends; the empty code (no
// code) counts as known.
func (c RejectCode) Known() bool {
	switch c {
	case "", RejectCodeInvalidCreditorAccountNumber:
		return true
	}
	return false
}

func (b *RejectBody) TAPType() string { return TypeReject }

// NewRejectMessage creates a new DIDComm message with a Reject body.
func NewRejectMessage(from string, to []string, thid string, body *RejectBody) (*didcomm.Message, error) {
	if !body.Code.Known() {
		return nil, fmt.Errorf("%w: unknown reject code %q", ErrInvalidBody, body.Code)
	}
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
