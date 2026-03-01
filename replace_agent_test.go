package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewReplaceAgentMessage(t *testing.T) {
	body := &ReplaceAgentBody{
		Original:    "did:pkh:eip155:1:0xabcda96D359eC26a11e2C2b3d8f8B8942d5Bfcdb",
		Replacement: &Agent{ID: "did:pkh:eip155:1:0x1234a96D359eC26a11e2C2b3d8f8B8942d5Bfcdb", Role: "settlementAddress"},
	}
	msg, err := NewReplaceAgentMessage("did:web:beneficiary.vasp", []string{"did:web:originator.vasp"}, "thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeReplaceAgent {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewReplaceAgentMessage_MissingFields(t *testing.T) {
	tests := []struct {
		name string
		body *ReplaceAgentBody
	}{
		{"missing original", &ReplaceAgentBody{Replacement: &Agent{ID: "did:web:new"}}},
		{"missing replacement", &ReplaceAgentBody{Original: "did:web:old"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewReplaceAgentMessage("from", nil, "thid", tt.body)
			if !errors.Is(err, ErrInvalidBody) {
				t.Errorf("expected ErrInvalidBody, got %v", err)
			}
		})
	}
}

func TestReplaceAgentBody_JSONRoundTrip(t *testing.T) {
	body := ReplaceAgentBody{
		Context:     TAPContext,
		Type:        TypeReplaceAgent,
		Original:    "did:web:old-agent",
		Replacement: &Agent{ID: "did:web:new-agent", Role: "settlementAddress"},
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got ReplaceAgentBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Original != body.Original || got.Replacement.ID != "did:web:new-agent" {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestReplaceAgent_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/replace-agent/valid.json")
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

	var body ReplaceAgentBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.Original != "did:pkh:eip155:1:0xabcda96D359eC26a11e2C2b3d8f8B8942d5Bfcdb" {
		t.Errorf("Original: got %q", body.Original)
	}
	if body.Replacement == nil || body.Replacement.Role != "settlementAddress" {
		t.Errorf("Replacement: got %+v", body.Replacement)
	}
}

func TestReplaceAgentBody_ParseBody(t *testing.T) {
	body := &ReplaceAgentBody{Original: "did:web:old", Replacement: &Agent{ID: "did:web:new"}}
	msg, err := NewReplaceAgentMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rab, ok := parsed.(*ReplaceAgentBody)
	if !ok {
		t.Fatalf("expected *ReplaceAgentBody, got %T", parsed)
	}
	if rab.Original != "did:web:old" {
		t.Errorf("Original: got %q", rab.Original)
	}
}
