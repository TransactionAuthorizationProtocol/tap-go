package tap

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestNewConfirmRelationshipMessage(t *testing.T) {
	body := &ConfirmRelationshipBody{
		Relationship: &Relationship{
			Type: "customer",
			Parties: []Party{
				{ID: "did:eg:alice"},
				{ID: "did:eg:bob"},
			},
		},
		Status: "confirmed",
	}
	msg, err := NewConfirmRelationshipMessage("did:web:beneficiary.vasp", []string{"did:web:originator.vasp"}, "thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeConfirmRelationship {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewConfirmRelationshipMessage_MissingFields(t *testing.T) {
	tests := []struct {
		name string
		body *ConfirmRelationshipBody
	}{
		{"missing relationship", &ConfirmRelationshipBody{Status: "confirmed"}},
		{"missing status", &ConfirmRelationshipBody{Relationship: &Relationship{Type: "customer", Parties: []Party{{ID: "a"}, {ID: "b"}}}}},
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

func TestConfirmRelationshipBody_JSONRoundTrip(t *testing.T) {
	body := ConfirmRelationshipBody{
		Context: TAPContext,
		Type:    TypeConfirmRelationship,
		Relationship: &Relationship{
			Type:    "customer",
			Parties: []Party{{ID: "did:eg:alice"}, {ID: "did:eg:bob"}},
		},
		Status:    "confirmed",
		ValidFrom: "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got ConfirmRelationshipBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Status != "confirmed" || got.Relationship.Type != "customer" {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestConfirmRelationshipBody_ParseBody(t *testing.T) {
	body := &ConfirmRelationshipBody{
		Relationship: &Relationship{Type: "customer", Parties: []Party{{ID: "a"}, {ID: "b"}}},
		Status:       "pending",
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
	if crb.Status != "pending" {
		t.Errorf("Status: got %q", crb.Status)
	}
}
