package tap

import (
	"encoding/json"
	"errors"
	"slices"
	"testing"
)

func TestNewConfirmRelationshipMessage(t *testing.T) {
	body := &ConfirmRelationshipBody{
		ID:   "did:pkh:eip155:1:0x1234a96D359eC26a11e2C2b3d8f8B8942d5Bfcdb",
		For:  NewForField("did:web:beneficiary.vasp"),
		Role: "SettlementAddress",
	}
	msg, err := NewConfirmRelationshipMessage("did:web:beneficiary.vasp", []string{"did:web:originator.vasp"}, "thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeConfirmRelationship {
		t.Errorf("Type: got %q", msg.Type)
	}
	if body.Type != TypeAgent {
		t.Errorf("body @type: got %q, want %q", body.Type, TypeAgent)
	}
}

func TestNewConfirmRelationshipMessage_MissingFields(t *testing.T) {
	tests := []struct {
		name string
		body *ConfirmRelationshipBody
	}{
		{"missing @id", &ConfirmRelationshipBody{For: NewForField("did:web:beneficiary.vasp")}},
		{"missing for", &ConfirmRelationshipBody{ID: "did:pkh:eip155:1:0x1234"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfirmRelationshipMessage("from", nil, "thid", tt.body)
			if !errors.Is(err, ErrInvalidBody) {
				t.Errorf("expected ErrInvalidBody, got %v", err)
			}
		})
	}
}

// TestConfirmRelationshipBody_SpecCanonical verifies the spec's canonical
// settlement-address example (TAIP-9 test case 1) round-trips with the exact
// JSON field names and casing.
func TestConfirmRelationshipBody_SpecCanonical(t *testing.T) {
	raw := []byte(`{
		"@context":"https://tap.rsvp/schema/1.0",
		"@type":"https://tap.rsvp/schema/1.0#Agent",
		"@id":"did:pkh:eip155:1:0x1234a96D359eC26a11e2C2b3d8f8B8942d5Bfcdb",
		"for":"did:web:beneficiary.vasp",
		"role":"SettlementAddress"
	}`)

	var got ConfirmRelationshipBody
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ID != "did:pkh:eip155:1:0x1234a96D359eC26a11e2C2b3d8f8B8942d5Bfcdb" {
		t.Errorf("@id: got %q", got.ID)
	}
	if got.For.String() != "did:web:beneficiary.vasp" {
		t.Errorf("for: got %q", got.For.String())
	}
	if got.Role != "SettlementAddress" {
		t.Errorf("role: got %q", got.Role)
	}
	if got.Type != TypeAgent {
		t.Errorf("@type: got %q, want %q", got.Type, TypeAgent)
	}

	// Re-marshal and confirm the spec JSON field names are present.
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatalf("unmarshal fields: %v", err)
	}
	for _, key := range []string{"@context", "@type", "@id", "for", "role"} {
		if _, ok := fields[key]; !ok {
			t.Errorf("marshalled body missing %q field; got keys %v", key, keysOf(fields))
		}
	}
}

// TestConfirmRelationshipBody_SpecCanonicalNoRole verifies the spec's second
// test case (VASP confirming it acts for its customer), where role is omitted.
func TestConfirmRelationshipBody_SpecCanonicalNoRole(t *testing.T) {
	raw := []byte(`{
		"@context":"https://tap.rsvp/schema/1.0",
		"@type":"https://tap.rsvp/schema/1.0#Agent",
		"@id":"did:web:beneficiary.vasp",
		"for":"did:eg:bob"
	}`)

	var got ConfirmRelationshipBody
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ID != "did:web:beneficiary.vasp" || got.For.String() != "did:eg:bob" {
		t.Errorf("mismatch: %+v", got)
	}

	data, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatalf("unmarshal fields: %v", err)
	}
	if _, ok := fields["role"]; ok {
		t.Errorf("role should be omitted when empty; got keys %v", keysOf(fields))
	}
}

func TestConfirmRelationshipBody_ForRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		for_ ForField
		want []string
	}{
		{"single owner", NewForField("did:web:beneficiary.vasp"), []string{"did:web:beneficiary.vasp"}},
		{"multiple owners", NewForField("did:eg:alice", "did:eg:bob"), []string{"did:eg:alice", "did:eg:bob"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := ConfirmRelationshipBody{
				Context: TAPContext,
				Type:    TypeAgent,
				ID:      "did:pkh:eip155:1:0x1234",
				For:     tt.for_,
			}

			data, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			var got ConfirmRelationshipBody
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if gotVals := got.For.Values(); !slices.Equal(gotVals, tt.want) {
				t.Errorf("For: got %v, want %v", gotVals, tt.want)
			}
		})
	}
}

func TestConfirmRelationshipBody_ParseBody(t *testing.T) {
	body := &ConfirmRelationshipBody{
		ID:   "did:pkh:eip155:1:0x1234",
		For:  NewForField("did:web:beneficiary.vasp"),
		Role: "SettlementAddress",
	}
	msg, err := NewConfirmRelationshipMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	crb, ok := parsed.(*ConfirmRelationshipBody)
	if !ok {
		t.Fatalf("expected *ConfirmRelationshipBody, got %T", parsed)
	}
	if crb.ID != "did:pkh:eip155:1:0x1234" {
		t.Errorf("@id: got %q", crb.ID)
	}
	if crb.For.String() != "did:web:beneficiary.vasp" {
		t.Errorf("for: got %q", crb.For.String())
	}
	if crb.Role != "SettlementAddress" {
		t.Errorf("role: got %q", crb.Role)
	}
}

func keysOf(m map[string]json.RawMessage) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
