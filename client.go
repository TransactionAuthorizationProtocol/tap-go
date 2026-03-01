package tap

import (
	"context"

	didcomm "github.com/Notabene-id/go-didcomm"
)

// TAPResult contains a typed TAP message after unpacking and parsing.
type TAPResult struct {
	Message *didcomm.Message
	Body    TAPBody

	// DIDComm envelope metadata
	Encrypted bool
	Signed    bool
	Anonymous bool
}

// Client wraps a didcomm.Client and adds TAP-specific typed message parsing.
type Client struct {
	DIDComm *didcomm.Client
}

// NewClient creates a new TAP client wrapping the given DIDComm client.
func NewClient(dc *didcomm.Client) *Client {
	return &Client{DIDComm: dc}
}

// Receive unpacks a DIDComm envelope and parses the TAP body into a typed struct.
func (c *Client) Receive(ctx context.Context, envelope []byte) (*TAPResult, error) {
	result, err := c.DIDComm.Unpack(ctx, envelope)
	if err != nil {
		return nil, err
	}

	body, err := ParseBody(result.Message)
	if err != nil {
		return nil, err
	}

	return &TAPResult{
		Message:   result.Message,
		Body:      body,
		Encrypted: result.Encrypted,
		Signed:    result.Signed,
		Anonymous: result.Anonymous,
	}, nil
}
