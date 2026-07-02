package tap

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestNewConnectMessage(t *testing.T) {
	body := &ConnectBody{
		Requester: &Party{ID: "did:web:b2b-service.example"},
		Principal: &Party{ID: "did:web:business-customer.example"},
		Agents: []Agent{
			{ID: "did:web:b2b-service.example", For: NewForField("did:web:b2b-service.example")},
		},
		Constraints: &TransactionConstraints{
			Purposes: []string{"BEXP"},
			Limits:   &Limits{PerTransaction: "10000.00", Currency: "USD"},
		},
	}
	msg, err := NewConnectMessage("did:web:b2b-service.example", []string{"did:web:vasp.example"}, body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeConnect {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewConnectMessage_MissingFields(t *testing.T) {
	base := ConnectBody{
		Requester:   &Party{ID: "did:web:req"},
		Principal:   &Party{ID: "did:web:princ"},
		Agents:      []Agent{{ID: "did:web:agent"}},
		Constraints: &TransactionConstraints{},
	}

	tests := []struct {
		name   string
		modify func(*ConnectBody)
	}{
		{"missing requester", func(b *ConnectBody) { b.Requester = nil }},
		{"missing principal", func(b *ConnectBody) { b.Principal = nil }},
		{"missing agents", func(b *ConnectBody) { b.Agents = nil }},
		{"missing constraints", func(b *ConnectBody) { b.Constraints = nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := base
			tt.modify(&b)
			_, err := NewConnectMessage("from", nil, &b)
			if !errors.Is(err, ErrInvalidBody) {
				t.Errorf("expected ErrInvalidBody, got %v", err)
			}
		})
	}
}

func TestConnectBody_JSONRoundTrip(t *testing.T) {
	body := ConnectBody{
		Context:   TAPContext,
		Type:      TypeConnect,
		Requester: &Party{ID: "did:web:req"},
		Principal: &Party{ID: "did:web:princ"},
		Agents:    []Agent{{ID: "did:web:agent"}},
		Constraints: &TransactionConstraints{
			Purposes: []string{"BEXP"},
			Limits:   &Limits{PerTransaction: "10000.00", Currency: "USD"},
		},
		Agreement: "https://example.com/terms",
		Expiry:    "2024-03-22T15:00:00Z",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got ConnectBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Requester.ID != "did:web:req" || got.Agreement != "https://example.com/terms" {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestConnect_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/connect/valid-b2b-connect.json")
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

	var body ConnectBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.Requester == nil || body.Requester.ID != "did:web:b2b-service.example" {
		t.Errorf("Requester: got %+v", body.Requester)
	}
	if body.Constraints == nil || len(body.Constraints.Purposes) != 2 {
		t.Errorf("Constraints.Purposes: got %v", body.Constraints)
	}
	if body.Constraints.Limits == nil || body.Constraints.Limits.PerTransaction != "10000.00" {
		t.Errorf("Constraints.Limits: got %v", body.Constraints.Limits)
	}
}

func TestConnectBody_ParseBody(t *testing.T) {
	body := &ConnectBody{
		Requester:   &Party{ID: "did:web:req"},
		Principal:   &Party{ID: "did:web:princ"},
		Agents:      []Agent{{ID: "did:web:agent"}},
		Constraints: &TransactionConstraints{Purposes: []string{"BEXP"}},
	}
	msg, err := NewConnectMessage("did:web:req", []string{"did:web:vasp"}, body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cb, ok := parsed.(*ConnectBody)
	if !ok {
		t.Fatalf("expected *ConnectBody, got %T", parsed)
	}
	if cb.Requester.ID != "did:web:req" {
		t.Errorf("Requester.ID: got %q", cb.Requester.ID)
	}
}

func TestNewConnectMessage_TrustConnection(t *testing.T) {
	body := &ConnectBody{
		ConnectionTypes: []string{ConnectionTypeDDQAccess},
		Action:          ConnectActionEstablish,
	}
	msg, err := NewConnectMessage("did:web:req.example", []string{"did:web:owner.example"}, body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(msg.Body, &raw); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	for _, field := range []string{"requester", "principal", "agents", "constraints"} {
		if _, ok := raw[field]; ok {
			t.Errorf("%s must be omitted on a trust connection", field)
		}
	}
	if string(raw["connectionTypes"]) != `["ddq-access"]` {
		t.Errorf("connectionTypes: got %s", raw["connectionTypes"])
	}
	if string(raw["action"]) != `"establish"` {
		t.Errorf("action: got %s", raw["action"])
	}
}

func TestNewConnectMessage_TransactionalTypeRequiresFields(t *testing.T) {
	body := &ConnectBody{
		ConnectionTypes: []string{ConnectionTypeTransaction},
	}
	_, err := NewConnectMessage("from", nil, body)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestConnectBody_JSONRoundTrip_TrustFields(t *testing.T) {
	body := ConnectBody{
		Context:         TAPContext,
		Type:            TypeConnect,
		ConnectionTypes: []string{ConnectionTypeDDQAccess, ConnectionTypeMutualTrust},
		Action:          ConnectActionUpdate,
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got ConnectBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.ConnectionTypes) != 2 || got.ConnectionTypes[0] != ConnectionTypeDDQAccess {
		t.Errorf("ConnectionTypes: got %v", got.ConnectionTypes)
	}
	if got.Action != ConnectActionUpdate {
		t.Errorf("Action: got %q", got.Action)
	}
}

func TestConnect_TestVectorEstablishDDQ(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/connect/valid-establish-ddq.json")
	if err != nil {
		t.Skipf("test vector not available: %v", err)
	}

	var tv struct {
		Body json.RawMessage `json:"body"`
	}
	if err := json.Unmarshal(data, &tv); err != nil {
		t.Fatalf("unmarshal test vector: %v", err)
	}

	var body ConnectBody
	if err := json.Unmarshal(tv.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if len(body.ConnectionTypes) == 0 || body.ConnectionTypes[0] != ConnectionTypeDDQAccess {
		t.Errorf("ConnectionTypes: got %v", body.ConnectionTypes)
	}
}
