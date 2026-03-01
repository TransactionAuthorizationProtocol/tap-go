package tap

import (
	"context"
	"encoding/json"
	"testing"

	didcomm "github.com/Notabene-id/go-didcomm"
)

func TestNewClient(t *testing.T) {
	dc := didcomm.NewClient(nil, nil)
	c := NewClient(dc)
	if c.DIDComm != dc {
		t.Error("DIDComm client not set correctly")
	}
}

func TestClient_ReceivePlainMessage(t *testing.T) {
	// Create a plain (unencrypted, unsigned) DIDComm message with a TAP body
	body := &TransferBody{
		Context: TAPContext,
		Type:    TypeTransfer,
		Asset:   "eip155:1/slip44:60",
		Amount:  "1.0",
		Agents:  []Agent{{ID: "did:web:originator.vasp", For: NewForField("did:eg:bob")}},
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	msg := &didcomm.Message{
		ID:   "test-123",
		Type: TypeTransfer,
		From: "did:web:originator.vasp",
		To:   []string{"did:web:beneficiary.vasp"},
		Body: rawBody,
	}

	envelope, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal message: %v", err)
	}

	// Use a client with nil resolvers (plain messages don't need them)
	dc := didcomm.NewClient(nil, nil)
	client := NewClient(dc)

	result, err := client.Receive(context.TODO(), envelope)
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}

	if result.Message.ID != "test-123" {
		t.Errorf("Message.ID: got %q", result.Message.ID)
	}
	if result.Encrypted || result.Signed || result.Anonymous {
		t.Error("expected plain message flags")
	}

	tb, ok := result.Body.(*TransferBody)
	if !ok {
		t.Fatalf("expected *TransferBody, got %T", result.Body)
	}
	if tb.Asset != "eip155:1/slip44:60" {
		t.Errorf("Asset: got %q", tb.Asset)
	}
}

func TestTAPResult_TypeAssertion(t *testing.T) {
	// Verify that we can type-assert the Body field to specific types
	body := &PaymentBody{
		Context:  TAPContext,
		Type:     TypePayment,
		Amount:   "100.00",
		Currency: "USD",
		Merchant: &Party{ID: "did:example:merchant"},
		Agents:   []Agent{{ID: "did:example:merchant"}},
	}
	rawBody, _ := json.Marshal(body)

	msg := &didcomm.Message{
		ID:   "test-pay",
		Type: TypePayment,
		Body: rawBody,
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("ParseBody: %v", err)
	}

	result := &TAPResult{
		Message: msg,
		Body:    parsed,
	}

	switch b := result.Body.(type) {
	case *PaymentBody:
		if b.Amount != "100.00" {
			t.Errorf("Amount: got %q", b.Amount)
		}
	default:
		t.Fatalf("unexpected type: %T", result.Body)
	}
}
