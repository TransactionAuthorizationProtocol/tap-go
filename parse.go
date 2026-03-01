package tap

import (
	"encoding/json"
	"fmt"

	didcomm "github.com/Notabene-id/go-didcomm"
)

// TAPBody is the interface implemented by all TAP message body types.
type TAPBody interface {
	TAPType() string
}

// IsTAPMessage returns true if the DIDComm message type is a known TAP type.
func IsTAPMessage(msg *didcomm.Message) bool {
	for _, t := range AllTypes() {
		if msg.Type == t {
			return true
		}
	}
	return false
}

// ParseBody unmarshals a DIDComm message body into the appropriate typed TAP body struct
// based on the message type.
func ParseBody(msg *didcomm.Message) (TAPBody, error) {
	var body TAPBody

	switch msg.Type {
	case TypeTransfer:
		body = &TransferBody{}
	case TypePayment:
		body = &PaymentBody{}
	case TypeExchange:
		body = &ExchangeBody{}
	case TypeQuote:
		body = &QuoteBody{}
	case TypeEscrow:
		body = &EscrowBody{}
	case TypeAuthorize:
		body = &AuthorizeBody{}
	case TypeAuthorizationRequired:
		body = &AuthorizationRequiredBody{}
	case TypeSettle:
		body = &SettleBody{}
	case TypeReject:
		body = &RejectBody{}
	case TypeCancel:
		body = &CancelBody{}
	case TypeRevert:
		body = &RevertBody{}
	case TypeCapture:
		body = &CaptureBody{}
	case TypeUpdateAgent:
		body = &UpdateAgentBody{}
	case TypeUpdateParty:
		body = &UpdatePartyBody{}
	case TypeAddAgents:
		body = &AddAgentsBody{}
	case TypeReplaceAgent:
		body = &ReplaceAgentBody{}
	case TypeRemoveAgent:
		body = &RemoveAgentBody{}
	case TypeConfirmRelationship:
		body = &ConfirmRelationshipBody{}
	case TypeUpdatePolicies:
		body = &UpdatePoliciesBody{}
	case TypeConnect:
		body = &ConnectBody{}
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownMessageType, msg.Type)
	}

	if err := json.Unmarshal(msg.Body, body); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidBody, err)
	}

	return body, nil
}
