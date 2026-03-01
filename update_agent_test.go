package tap

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestNewUpdateAgentMessage(t *testing.T) {
	body := &UpdateAgentBody{
		Agent:  &Agent{ID: "did:web:new.vasp", For: NewForField("did:eg:bob")},
		Reason: "Key rotation",
	}
	msg, err := NewUpdateAgentMessage("did:web:originator.vasp", []string{"did:web:beneficiary.vasp"}, "thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeUpdateAgent {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewUpdateAgentMessage_MissingAgent(t *testing.T) {
	body := &UpdateAgentBody{}
	_, err := NewUpdateAgentMessage("from", nil, "thid", body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestUpdateAgentBody_JSONRoundTrip(t *testing.T) {
	body := UpdateAgentBody{
		Context:     TAPContext,
		Type:        TypeUpdateAgent,
		Agent:       &Agent{ID: "did:web:agent", Role: "custodian"},
		PreviousDID: "did:web:old-agent",
		Reason:      "Migration",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got UpdateAgentBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Agent.ID != "did:web:agent" || got.PreviousDID != "did:web:old-agent" {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestUpdateAgentBody_ParseBody(t *testing.T) {
	body := &UpdateAgentBody{Agent: &Agent{ID: "did:web:new"}}
	msg, err := NewUpdateAgentMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	uab, ok := parsed.(*UpdateAgentBody)
	if !ok {
		t.Fatalf("expected *UpdateAgentBody, got %T", parsed)
	}
	if uab.Agent.ID != "did:web:new" {
		t.Errorf("Agent.ID: got %q", uab.Agent.ID)
	}
}
