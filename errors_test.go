package tap

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrInvalidBody", ErrInvalidBody},
		{"ErrUnknownMessageType", ErrUnknownMessageType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatal("error should not be nil")
			}
			if tt.err.Error() == "" {
				t.Fatal("error message should not be empty")
			}
		})
	}
}

func TestErrors_Wrapping(t *testing.T) {
	wrapped := fmt.Errorf("%w: missing field", ErrInvalidBody)
	if !errors.Is(wrapped, ErrInvalidBody) {
		t.Error("wrapped error should match ErrInvalidBody")
	}

	wrapped2 := fmt.Errorf("%w: unknown", ErrUnknownMessageType)
	if !errors.Is(wrapped2, ErrUnknownMessageType) {
		t.Error("wrapped error should match ErrUnknownMessageType")
	}
}
