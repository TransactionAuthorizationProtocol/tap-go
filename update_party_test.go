package tap

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestNewUpdatePartyMessage(t *testing.T) {
	body := &UpdatePartyBody{
		Party: &Party{ID: "did:eg:bob", Name: "Bob Updated"},
		Role:  "originator",
	}
	msg, err := NewUpdatePartyMessage("did:web:originator.vasp", []string{"did:web:beneficiary.vasp"}, "thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeUpdateParty {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewUpdatePartyMessage_MissingFields(t *testing.T) {
	tests := []struct {
		name string
		body *UpdatePartyBody
	}{
		{"missing party", &UpdatePartyBody{Role: "originator"}},
		{"missing role", &UpdatePartyBody{Party: &Party{ID: "did:eg:bob"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUpdatePartyMessage("from", nil, "thid", tt.body)
			if !errors.Is(err, ErrInvalidBody) {
				t.Errorf("expected ErrInvalidBody, got %v", err)
			}
		})
	}
}

func TestUpdatePartyBody_JSONRoundTrip(t *testing.T) {
	body := UpdatePartyBody{
		Context: TAPContext,
		Type:    TypeUpdateParty,
		Party:   &Party{ID: "did:eg:alice", Name: "Alice"},
		Role:    "beneficiary",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got UpdatePartyBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Party.Name != "Alice" || got.Role != "beneficiary" {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestUpdatePartyBody_ParseBody(t *testing.T) {
	body := &UpdatePartyBody{Party: &Party{ID: "did:eg:alice"}, Role: "beneficiary"}
	msg, err := NewUpdatePartyMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	upb, ok := parsed.(*UpdatePartyBody)
	if !ok {
		t.Fatalf("expected *UpdatePartyBody, got %T", parsed)
	}
	if upb.Role != "beneficiary" {
		t.Errorf("Role: got %q", upb.Role)
	}
}
