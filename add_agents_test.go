package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewAddAgentsMessage(t *testing.T) {
	body := &AddAgentsBody{
		Agents: []Agent{
			{ID: "did:web:beneficiary.vasp", For: NewForField("did:eg:alice")},
		},
	}
	msg, err := NewAddAgentsMessage("did:web:beneficiary.vasp", []string{"did:web:originator.vasp"}, "thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeAddAgents {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewAddAgentsMessage_MissingAgents(t *testing.T) {
	body := &AddAgentsBody{}
	_, err := NewAddAgentsMessage("from", nil, "thid", body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestAddAgentsBody_JSONRoundTrip(t *testing.T) {
	body := AddAgentsBody{
		Context: TAPContext,
		Type:    TypeAddAgents,
		Agents: []Agent{
			{ID: "did:web:beneficiary.vasp", For: NewForField("did:eg:alice")},
			{ID: "did:web:wallet.service", For: NewForField("did:web:beneficiary.vasp"), Role: "custodian"},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got AddAgentsBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.Agents) != 2 {
		t.Errorf("Agents: got %d, want 2", len(got.Agents))
	}
}

func TestAddAgents_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/add-agents/valid.json")
	if err != nil {
		t.Skipf("test vector not available: %v", err)
	}

	var tv struct {
		Message struct {
			Body json.RawMessage `json:"body"`
		} `json:"message"`
	}
	if err := json.Unmarshal(data, &tv); err != nil {
		t.Fatalf("unmarshal test vector: %v", err)
	}

	var body AddAgentsBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if len(body.Agents) != 3 {
		t.Errorf("Agents: got %d, want 3", len(body.Agents))
	}
}

func TestAddAgentsBody_ParseBody(t *testing.T) {
	body := &AddAgentsBody{Agents: []Agent{{ID: "did:web:new"}}}
	msg, err := NewAddAgentsMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	aab, ok := parsed.(*AddAgentsBody)
	if !ok {
		t.Fatalf("expected *AddAgentsBody, got %T", parsed)
	}
	if len(aab.Agents) != 1 {
		t.Errorf("Agents: got %d", len(aab.Agents))
	}
}
