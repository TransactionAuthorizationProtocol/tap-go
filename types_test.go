package tap

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestParty_JSONRoundTrip(t *testing.T) {
	party := Party{
		ID:   "did:eg:bob",
		Type: "Party",
		Name: "Bob",
	}

	data, err := json.Marshal(party)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got Party
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ID != party.ID || got.Type != party.Type || got.Name != party.Name {
		t.Errorf("round-trip mismatch: got %+v, want %+v", got, party)
	}
}

func TestAgent_JSONRoundTrip(t *testing.T) {
	agent := Agent{
		ID:   "did:web:originator.vasp",
		For:  NewForField("did:eg:bob"),
		Role: "SettlementAddress",
		Policies: []Policy{
			{Type: "RequireAuthorization", FromAgent: "originator"},
		},
	}

	data, err := json.Marshal(agent)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got Agent
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ID != agent.ID {
		t.Errorf("ID: got %q, want %q", got.ID, agent.ID)
	}
	if got.For.String() != "did:eg:bob" {
		t.Errorf("For: got %q, want %q", got.For.String(), "did:eg:bob")
	}
	if got.Role != "SettlementAddress" {
		t.Errorf("Role: got %q, want %q", got.Role, "SettlementAddress")
	}
	if len(got.Policies) != 1 {
		t.Fatalf("Policies: got %d, want 1", len(got.Policies))
	}
	if got.Policies[0].Type != "RequireAuthorization" {
		t.Errorf("Policies[0].Type: got %q, want %q", got.Policies[0].Type, "RequireAuthorization")
	}
}

func TestForField_SingleDID(t *testing.T) {
	f := NewForField("did:eg:alice")
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"did:eg:alice"` {
		t.Errorf("marshal single: got %s, want %q", data, "did:eg:alice")
	}

	var got ForField
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.String() != "did:eg:alice" {
		t.Errorf("unmarshal single: got %q, want %q", got.String(), "did:eg:alice")
	}
}

func TestForField_MultipleDIDs(t *testing.T) {
	f := NewForField("did:eg:alice", "did:eg:bob")
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `["did:eg:alice","did:eg:bob"]` {
		t.Errorf("marshal multi: got %s", data)
	}

	var got ForField
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	vals := got.Values()
	if len(vals) != 2 || vals[0] != "did:eg:alice" || vals[1] != "did:eg:bob" {
		t.Errorf("unmarshal multi: got %v", vals)
	}
}

func TestForField_Empty(t *testing.T) {
	var f ForField
	if !f.IsEmpty() {
		t.Error("expected empty ForField")
	}
	if f.String() != "" {
		t.Errorf("expected empty string, got %q", f.String())
	}
}

func TestPolicy_JSONRoundTrip(t *testing.T) {
	policy := Policy{
		Type:                   "RequirePresentation",
		FromAgent:              "beneficiary",
		AboutParty:             "beneficiary",
		Purpose:                "FATF Travel Rule Compliance",
		PresentationDefinition: "https://tap.rsvp/presentation-definitions/ivms-101/eu/tfr",
	}

	data, err := json.Marshal(policy)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got Policy
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Type != policy.Type || got.FromAgent != policy.FromAgent {
		t.Errorf("round-trip mismatch: got %+v, want %+v", got, policy)
	}
}

func TestPolicy_FromFieldRoundTrip(t *testing.T) {
	in := Policy{Type: "RequireAuthorization", From: "did:web:vasp.example"}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !bytes.Contains(data, []byte(`"from":"did:web:vasp.example"`)) {
		t.Fatalf("missing from in JSON: %s", string(data))
	}
	var out Policy
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.From != in.From {
		t.Errorf("From = %q, want %q", out.From, in.From)
	}
}

func TestTransactionConstraints_JSONRoundTrip(t *testing.T) {
	tc := TransactionConstraints{
		Purposes:         []string{"BEXP", "SUPP"},
		CategoryPurposes: []string{"CASH"},
		Limits: &Limits{
			PerTransaction: "10000.00",
			PerDay:         "50000.00",
			Currency:       "USD",
		},
		AllowedBeneficiaries: []Party{
			{ID: "did:example:vendor-1", Name: "Approved Vendor 1"},
		},
		AllowedSettlementAddresses: []string{"eip155:1:0x742d35Cc6e4dfE2eDFaD2C0b91A8b0780EDAEb58"},
		AllowedAssets:              []string{"eip155:1/slip44:60"},
	}

	data, err := json.Marshal(tc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got TransactionConstraints
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(got.Purposes) != 2 || got.Purposes[0] != "BEXP" {
		t.Errorf("Purposes mismatch: got %v", got.Purposes)
	}
	if got.Limits == nil || got.Limits.PerTransaction != "10000.00" {
		t.Errorf("Limits mismatch: got %v", got.Limits)
	}
}

func TestInvoice_JSONRoundTrip(t *testing.T) {
	inv := Invoice{
		ID:           "INV001",
		IssueDate:    "2025-04-20",
		CurrencyCode: "USD",
		LineItems: []LineItem{
			{
				ID:          "1",
				Description: "Widget",
				Quantity:    2,
				UnitPrice:   10.00,
				LineTotal:   20.00,
			},
		},
		Total: 20.00,
	}

	data, err := json.Marshal(inv)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got Invoice
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ID != inv.ID || got.Total != inv.Total {
		t.Errorf("round-trip mismatch: got %+v", got)
	}
	if len(got.LineItems) != 1 || got.LineItems[0].Description != "Widget" {
		t.Errorf("LineItems mismatch: got %v", got.LineItems)
	}
}

func TestRelationship_JSONRoundTrip(t *testing.T) {
	rel := Relationship{
		Type: "customer",
		Parties: []Party{
			{ID: "did:eg:alice"},
			{ID: "did:eg:bob"},
		},
	}

	data, err := json.Marshal(rel)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got Relationship
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Type != "customer" || len(got.Parties) != 2 {
		t.Errorf("round-trip mismatch: got %+v", got)
	}
}
