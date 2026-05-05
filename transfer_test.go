package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewTransferMessage(t *testing.T) {
	body := &TransferBody{
		Asset:      "eip155:1/slip44:60",
		Amount:     "1.23",
		Originator: &Party{ID: "did:eg:bob"},
		Agents: []Agent{
			{ID: "did:web:originator.vasp", For: NewForField("did:eg:bob")},
		},
	}

	msg, err := NewTransferMessage("did:web:originator.vasp", []string{"did:web:beneficiary.vasp"}, body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeTransfer {
		t.Errorf("Type: got %q, want %q", msg.Type, TypeTransfer)
	}
	if msg.From != "did:web:originator.vasp" {
		t.Errorf("From: got %q", msg.From)
	}
	if msg.ID == "" {
		t.Error("ID should not be empty")
	}

	// Verify body was marshaled correctly
	var got TransferBody
	if err := json.Unmarshal(msg.Body, &got); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if got.Context != TAPContext {
		t.Errorf("Context: got %q, want %q", got.Context, TAPContext)
	}
	if got.Type != TypeTransfer {
		t.Errorf("Body Type: got %q, want %q", got.Type, TypeTransfer)
	}
	if got.Asset != "eip155:1/slip44:60" {
		t.Errorf("Asset: got %q", got.Asset)
	}
}

func TestNewTransferMessage_MissingAsset(t *testing.T) {
	body := &TransferBody{
		Agents: []Agent{{ID: "did:web:originator.vasp"}},
	}
	_, err := NewTransferMessage("did:web:originator.vasp", []string{"did:web:beneficiary.vasp"}, body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestNewTransferMessage_MissingAgents(t *testing.T) {
	body := &TransferBody{
		Asset: "eip155:1/slip44:60",
	}
	_, err := NewTransferMessage("did:web:originator.vasp", []string{"did:web:beneficiary.vasp"}, body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestTransferBody_JSONRoundTrip(t *testing.T) {
	body := TransferBody{
		Context:     TAPContext,
		Type:        TypeTransfer,
		Asset:       "eip155:1/slip44:60",
		Amount:      "1.23",
		Originator:  &Party{ID: "did:eg:bob"},
		Beneficiary: &Party{ID: "did:eg:alice"},
		Agents: []Agent{
			{ID: "did:web:originator.vasp", For: NewForField("did:eg:bob")},
			{ID: "did:web:beneficiary.vasp", For: NewForField("did:eg:alice")},
		},
		SettlementID: "eip155:1:tx/0x3edb98c24d46d148eb926c714f4fbaa117c47b0c0821f38bfce9763604457c33",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got TransferBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Asset != body.Asset || got.Amount != body.Amount || got.SettlementID != body.SettlementID {
		t.Errorf("round-trip mismatch: got %+v", got)
	}
}

func TestTransferBody_ParseBody(t *testing.T) {
	body := &TransferBody{
		Asset:  "eip155:1/slip44:60",
		Amount: "1.23",
		Agents: []Agent{{ID: "did:web:originator.vasp", For: NewForField("did:eg:bob")}},
	}

	msg, err := NewTransferMessage("did:web:originator.vasp", []string{"did:web:beneficiary.vasp"}, body)
	if err != nil {
		t.Fatalf("create message: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse body: %v", err)
	}

	tb, ok := parsed.(*TransferBody)
	if !ok {
		t.Fatalf("expected *TransferBody, got %T", parsed)
	}
	if tb.Asset != "eip155:1/slip44:60" {
		t.Errorf("Asset: got %q", tb.Asset)
	}
}

func TestTransferBody_TransactionValue(t *testing.T) {
	body := &TransferBody{
		Asset:  "eip155:1/erc20:0x1234567890abcdef1234567890abcdef12345678",
		Amount: "500",
		TransactionValue: &TransactionValue{
			Amount:   "1250.00",
			Currency: "EUR",
		},
		Agents: []Agent{{ID: "did:web:originator.vasp", For: NewForField("did:eg:bob")}},
	}

	msg, err := NewTransferMessage("did:web:originator.vasp", []string{"did:web:beneficiary.vasp"}, body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	tb := parsed.(*TransferBody)
	if tb.TransactionValue == nil {
		t.Fatal("TransactionValue: got nil")
	}
	if tb.TransactionValue.Amount != "1250.00" {
		t.Errorf("TransactionValue.Amount: got %q", tb.TransactionValue.Amount)
	}
	if tb.TransactionValue.Currency != "EUR" {
		t.Errorf("TransactionValue.Currency: got %q", tb.TransactionValue.Currency)
	}

	// transactionValue is optional — confirm omitempty when nil
	bodyNoValue := &TransferBody{
		Asset:  "eip155:1/slip44:60",
		Agents: []Agent{{ID: "did:web:originator.vasp"}},
	}
	msg2, err := NewTransferMessage("did:web:originator.vasp", nil, bodyNoValue)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if string(msg2.Body) == "" || jsonContains(msg2.Body, "transactionValue") {
		t.Errorf("expected no transactionValue field, body=%s", msg2.Body)
	}
}

func jsonContains(body []byte, key string) bool {
	for i := 0; i < len(body)-len(key); i++ {
		if string(body[i:i+len(key)]) == key {
			return true
		}
	}
	return false
}

func TestTransfer_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/transfer/valid.json")
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

	var body TransferBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.Asset != "eip155:1/slip44:60" {
		t.Errorf("Asset: got %q", body.Asset)
	}
	if body.Amount != "1.23" {
		t.Errorf("Amount: got %q", body.Amount)
	}
	if len(body.Agents) != 3 {
		t.Errorf("Agents: got %d, want 3", len(body.Agents))
	}
	if body.Context != TAPContext {
		t.Errorf("Context: got %q", body.Context)
	}
}

func TestTransfer_TestVectorMinimal(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/transfer/minimal.json")
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

	var body TransferBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.Asset != "eip155:1/slip44:60" {
		t.Errorf("Asset: got %q", body.Asset)
	}
	if len(body.Agents) != 1 {
		t.Errorf("Agents: got %d, want 1", len(body.Agents))
	}
}
