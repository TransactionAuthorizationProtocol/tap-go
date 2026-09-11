package tap

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	didcomm "github.com/notabene-id/go-didcomm"
)

// UpdatePartyBody represents the body of a TAP UpdateParty message (TAIP-6).
type UpdatePartyBody struct {
	Context       string `json:"@context"`
	Type          string `json:"@type"`
	Party         *Party `json:"party"`
	PartyType     string `json:"partyType"`
	PreviousParty *Party `json:"previousParty,omitempty"`
	Reason        string `json:"reason,omitempty"`
	Effective     string `json:"effective,omitempty"`
}

// UnmarshalJSON reads partyType, falling back to the "role" key that earlier
// tap-go releases emitted in its place. Peers still on the old spelling keep
// working; everything this library sends uses partyType.
func (b *UpdatePartyBody) UnmarshalJSON(data []byte) error {
	type alias UpdatePartyBody
	aux := struct {
		*alias
		LegacyRole string `json:"role"`
	}{alias: (*alias)(b)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if b.PartyType == "" {
		b.PartyType = aux.LegacyRole
	}
	return nil
}

func (b *UpdatePartyBody) TAPType() string { return TypeUpdateParty }

// NewUpdatePartyMessage creates a new DIDComm message with an UpdateParty body.
func NewUpdatePartyMessage(from string, to []string, thid string, body *UpdatePartyBody) (*didcomm.Message, error) {
	if body.Party == nil {
		return nil, fmt.Errorf("%w: missing party", ErrInvalidBody)
	}
	if body.PartyType == "" {
		return nil, fmt.Errorf("%w: missing partyType", ErrInvalidBody)
	}

	body.Context = TAPContext
	body.Type = TypeUpdateParty

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeUpdateParty,
		From: from,
		To:   to,
		Thid: thid,
		Body: rawBody,
	}, nil
}
