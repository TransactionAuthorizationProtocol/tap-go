package tap

import (
	"context"

	didcomm "github.com/notabene-id/go-didcomm"
)

// TAPResult contains a typed TAP message after unpacking and parsing.
type TAPResult struct {
	Message *didcomm.Message
	Body    TAPBody

	// SenderDID is the cryptographically verified sender (equal to Message.From
	// for authenticated messages, empty when unverified). Encrypted reports
	// whether the envelope was encrypted; Anonymous reports the absence of a
	// verified sender.
	SenderDID string
	Encrypted bool
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

// Receive unpacks a DIDComm envelope and parses the TAP body into a typed
// struct. It authenticates the sender: plain and anonymously-encrypted messages
// are rejected. Use ReceiveUnverified to accept those.
func (c *Client) Receive(ctx context.Context, envelope []byte) (*TAPResult, error) {
	msg, meta, err := c.DIDComm.Unpack(ctx, envelope)
	if err != nil {
		return nil, err
	}
	return c.toResult(msg, meta)
}

// ReceiveUnverified unpacks without requiring an authenticated sender. The
// returned SenderDID may be empty, in which case Message.From is not verified
// and must not be trusted.
func (c *Client) ReceiveUnverified(ctx context.Context, envelope []byte) (*TAPResult, error) {
	msg, meta, err := c.DIDComm.UnpackUnverified(ctx, envelope)
	if err != nil {
		return nil, err
	}
	return c.toResult(msg, meta)
}

func (c *Client) toResult(msg *didcomm.Message, meta *didcomm.Metadata) (*TAPResult, error) {
	body, err := ParseBody(msg)
	if err != nil {
		return nil, err
	}
	return &TAPResult{
		Message:   msg,
		Body:      body,
		SenderDID: meta.SenderDID,
		Encrypted: meta.Encrypted,
		Anonymous: meta.SenderDID == "",
	}, nil
}
