package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewRemoveAgentMessage(t *testing.T) {
	body := &RemoveAgentBody{
		Agent: "did:pkh:eip155:1:0xabcda96D359eC26a11e2C2b3d8f8B8942d5Bfcdb",
	}
	msg, err := NewRemoveAgentMessage("did:web:beneficiary.vasp", []string{"did:web:originator.vasp"}, "thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeRemoveAgent {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewRemoveAgentMessage_MissingAgent(t *testing.T) {
	body := &RemoveAgentBody{}
	_, err := NewRemoveAgentMessage("from", nil, "thid", body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestRemoveAgentBody_JSONRoundTrip(t *testing.T) {
	body := RemoveAgentBody{
		Context: TAPContext,
		Type:    TypeRemoveAgent,
		Agent:   "did:web:agent-to-remove",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got RemoveAgentBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Agent != body.Agent {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestRemoveAgent_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/remove-agent/valid.json")
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

	var body RemoveAgentBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.Agent != "did:pkh:eip155:1:0xabcda96D359eC26a11e2C2b3d8f8B8942d5Bfcdb" {
		t.Errorf("Agent: got %q", body.Agent)
	}
}

func TestRemoveAgentBody_ParseBody(t *testing.T) {
	body := &RemoveAgentBody{Agent: "did:web:old"}
	msg, err := NewRemoveAgentMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rab, ok := parsed.(*RemoveAgentBody)
	if !ok {
		t.Fatalf("expected *RemoveAgentBody, got %T", parsed)
	}
	if rab.Agent != "did:web:old" {
		t.Errorf("Agent: got %q", rab.Agent)
	}
}
