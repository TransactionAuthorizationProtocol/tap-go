package tap

import "errors"

var (
	// ErrInvalidBody is returned when a message body is missing required fields or malformed.
	ErrInvalidBody = errors.New("tap: invalid body")

	// ErrUnknownMessageType is returned when a message type is not recognized as a TAP type.
	ErrUnknownMessageType = errors.New("tap: unknown message type")
)
