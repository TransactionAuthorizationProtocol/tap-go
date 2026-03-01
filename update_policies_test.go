package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewUpdatePoliciesMessage(t *testing.T) {
	body := &UpdatePoliciesBody{
		Policies: []Policy{
			{Type: "RequireAuthorization", FromAgent: "beneficiary", Purpose: "FATF Travel Rule Compliance"},
		},
	}
	msg, err := NewUpdatePoliciesMessage("did:web:originator.vasp", []string{"did:web:beneficiary.vasp"}, "thread-1", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeUpdatePolicies {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewUpdatePoliciesMessage_MissingPolicies(t *testing.T) {
	body := &UpdatePoliciesBody{}
	_, err := NewUpdatePoliciesMessage("from", nil, "thid", body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestUpdatePoliciesBody_JSONRoundTrip(t *testing.T) {
	body := UpdatePoliciesBody{
		Context: TAPContext,
		Type:    TypeUpdatePolicies,
		Policies: []Policy{
			{Type: "RequirePresentation", FromAgent: "beneficiary", AboutParty: "beneficiary"},
			{Type: "RequireRelationshipConfirmation", FromRole: "SettlementAddress"},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got UpdatePoliciesBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.Policies) != 2 {
		t.Errorf("Policies: got %d, want 2", len(got.Policies))
	}
}

func TestUpdatePolicies_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/policy-management/valid-policies.json")
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

	var body UpdatePoliciesBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if len(body.Policies) != 3 {
		t.Errorf("Policies: got %d, want 3", len(body.Policies))
	}
}

func TestUpdatePoliciesBody_ParseBody(t *testing.T) {
	body := &UpdatePoliciesBody{
		Policies: []Policy{{Type: "RequireAuthorization"}},
	}
	msg, err := NewUpdatePoliciesMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	upb, ok := parsed.(*UpdatePoliciesBody)
	if !ok {
		t.Fatalf("expected *UpdatePoliciesBody, got %T", parsed)
	}
	if len(upb.Policies) != 1 {
		t.Errorf("Policies: got %d", len(upb.Policies))
	}
}
