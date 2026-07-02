package tap

import (
	"encoding/json"
	"fmt"
	"slices"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/google/uuid"
)

// Connection types carried in ConnectBody.ConnectionTypes (TAIP-15).
const (
	ConnectionTypeTransaction = "transaction"
	ConnectionTypeDDQAccess   = "ddq-access"
	ConnectionTypeMutualTrust = "mutual-trust"
	ConnectionTypeWhitelist   = "whitelist"
)

// Connection lifecycle actions carried in ConnectBody.Action (TAIP-15).
const (
	ConnectActionEstablish = "establish"
	ConnectActionUpdate    = "update"
)

// ConnectBody represents the body of a TAP Connect message (TAIP-15).
type ConnectBody struct {
	Context         string                  `json:"@context"`
	Type            string                  `json:"@type"`
	ConnectionTypes []string                `json:"connectionTypes,omitempty"`
	Action          string                  `json:"action,omitempty"`
	Requester       *Party                  `json:"requester,omitempty"`
	Principal       *Party                  `json:"principal,omitempty"`
	Agents          []Agent                 `json:"agents,omitempty"`
	Constraints     *TransactionConstraints `json:"constraints,omitempty"`
	Agreement       string                  `json:"agreement,omitempty"`
	Expiry          string                  `json:"expiry,omitempty"`
}

func (b *ConnectBody) TAPType() string { return TypeConnect }

// NewConnectMessage creates a new DIDComm message with a Connect body.
// Requester, principal, agents, and constraints are validated for
// transactional connections only.
func NewConnectMessage(from string, to []string, body *ConnectBody) (*didcomm.Message, error) {
	if transactional(body.ConnectionTypes) {
		if body.Requester == nil {
			return nil, fmt.Errorf("%w: missing requester", ErrInvalidBody)
		}
		if body.Principal == nil {
			return nil, fmt.Errorf("%w: missing principal", ErrInvalidBody)
		}
		if len(body.Agents) == 0 {
			return nil, fmt.Errorf("%w: missing agents", ErrInvalidBody)
		}
		if body.Constraints == nil {
			return nil, fmt.Errorf("%w: missing constraints", ErrInvalidBody)
		}
	}

	body.Context = TAPContext
	body.Type = TypeConnect

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	return &didcomm.Message{
		ID:   uuid.New().String(),
		Type: TypeConnect,
		From: from,
		To:   to,
		Body: rawBody,
	}, nil
}

// transactional reports whether the connection types describe a transactional
// connection. An empty list is a pre-revision Connect, which is always
// transactional.
func transactional(connectionTypes []string) bool {
	return len(connectionTypes) == 0 || slices.Contains(connectionTypes, ConnectionTypeTransaction)
}
