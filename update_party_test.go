package tap

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestNewUpdatePartyMessage(t *testing.T) {
	body := &UpdatePartyBody{
		Party:     &Party{ID: "did:eg:bob", Name: "Bob Updated"},
		PartyType: "originator",
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
		{"missing party", &UpdatePartyBody{PartyType: "originator"}},
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
		Context:   TAPContext,
		Type:      TypeUpdateParty,
		Party:     &Party{ID: "did:eg:alice", Name: "Alice"},
		PartyType: "beneficiary",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got UpdatePartyBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Party.Name != "Alice" || got.PartyType != "beneficiary" {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestUpdatePartyBody_ParseBody(t *testing.T) {
	body := &UpdatePartyBody{Party: &Party{ID: "did:eg:alice"}, PartyType: "beneficiary"}
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
	if upb.PartyType != "beneficiary" {
		t.Errorf("PartyType: got %q", upb.PartyType)
	}
}

// Earlier tap-go releases tagged this field "role". Peers still on that
// spelling must keep working while everything we send uses "partyType".
func TestUpdatePartyBody_UnmarshalAcceptsLegacyRoleKey(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"partyType", `{"partyType":"beneficiary"}`, "beneficiary"},
		{"legacy role", `{"role":"originator"}`, "originator"},
		{"partyType wins", `{"partyType":"beneficiary","role":"originator"}`, "beneficiary"},
		{"neither", `{}`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body UpdatePartyBody
			if err := json.Unmarshal([]byte(tt.raw), &body); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if body.PartyType != tt.want {
				t.Errorf("PartyType: got %q, want %q", body.PartyType, tt.want)
			}
		})
	}
}

func TestUpdatePartyBody_MarshalsPartyType(t *testing.T) {
	msg, err := NewUpdatePartyMessage("did:web:a", []string{"did:web:b"}, "thid", &UpdatePartyBody{
		Party:     &Party{ID: "did:eg:alice"},
		PartyType: "beneficiary",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(msg.Body, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := raw["role"]; ok {
		t.Error("body still emits the legacy role key")
	}
	if string(raw["partyType"]) != `"beneficiary"` {
		t.Errorf("partyType: got %s", raw["partyType"])
	}
}
