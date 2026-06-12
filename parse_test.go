package tap

import (
	"encoding/json"
	"errors"
	"testing"

	didcomm "github.com/Notabene-id/go-didcomm"
)

func TestIsTAPMessage(t *testing.T) {
	tests := []struct {
		msgType string
		want    bool
	}{
		{TypeTransfer, true},
		{TypePayment, true},
		{TypeRFQ, true},
		{TypeQuote, true},
		{TypeLock, true},
		{TypeAuthorize, true},
		{TypeAuthorizationRequired, true},
		{TypeSettle, true},
		{TypeReject, true},
		{TypeCancel, true},
		{TypeRevert, true},
		{TypeCapture, true},
		{TypeUpdateAgent, true},
		{TypeUpdateParty, true},
		{TypeAddAgents, true},
		{TypeReplaceAgent, true},
		{TypeRemoveAgent, true},
		{TypeConfirmRelationship, true},
		{TypeUpdatePolicies, true},
		{TypeConnect, true},
		{"https://example.com/unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		msg := &didcomm.Message{Type: tt.msgType}
		if got := IsTAPMessage(msg); got != tt.want {
			t.Errorf("IsTAPMessage(%q): got %v, want %v", tt.msgType, got, tt.want)
		}
	}
}

func TestParseBody_AllTypes(t *testing.T) {
	tests := []struct {
		name    string
		msgType string
		body    any
	}{
		{"Transfer", TypeTransfer, &TransferBody{Context: TAPContext, Type: TypeTransfer, Asset: "ETH", Agents: []Agent{{ID: "a"}}}},
		{"Payment", TypePayment, &PaymentBody{Context: TAPContext, Type: TypePayment, Amount: "100", Currency: "USD", Merchant: &Party{ID: "m"}, Agents: []Agent{{ID: "a"}}}},
		{"RFQ", TypeRFQ, &RFQBody{Context: TAPContext, Type: TypeRFQ, FromAssets: []string{"ETH"}, ToAssets: []string{"USD"}, FromAmount: "1", Requester: &Party{ID: "r"}, Agents: []Agent{{ID: "a"}}}},
		{"Quote", TypeQuote, &QuoteBody{Context: TAPContext, Type: TypeQuote, FromAsset: "ETH", ToAsset: "USD", FromAmount: "1", ToAmount: "3000", Provider: &Party{ID: "p"}, Agents: []Agent{{ID: "a"}}, ExpiresAt: "2025-01-01T00:00:00Z"}},
		{"Lock", TypeLock, &LockBody{Context: TAPContext, Type: TypeLock, Asset: "ETH", Amount: "100", Originator: &Party{ID: "o"}, Beneficiary: &Party{ID: "b"}, Expiry: "2025-01-01T00:00:00Z", Agents: []Agent{{ID: "a"}}}},
		{"Authorize", TypeAuthorize, &AuthorizeBody{Context: TAPContext, Type: TypeAuthorize}},
		{"AuthorizationRequired", TypeAuthorizationRequired, &AuthorizationRequiredBody{Context: TAPContext, Type: TypeAuthorizationRequired, AuthorizationURL: "https://example.com", Expires: "2025-01-01T00:00:00Z"}},
		{"Settle", TypeSettle, &SettleBody{Context: TAPContext, Type: TypeSettle, SettlementAddress: "eip155:1:0x1234"}},
		{"Reject", TypeReject, &RejectBody{Context: TAPContext, Type: TypeReject}},
		{"Cancel", TypeCancel, &CancelBody{Context: TAPContext, Type: TypeCancel, By: "originator"}},
		{"Revert", TypeRevert, &RevertBody{Context: TAPContext, Type: TypeRevert, SettlementAddress: "eip155:1:0x1234", Reason: "test"}},
		{"Capture", TypeCapture, &CaptureBody{Context: TAPContext, Type: TypeCapture}},
		{"UpdateAgent", TypeUpdateAgent, &UpdateAgentBody{Context: TAPContext, Type: TypeUpdateAgent, Agent: &Agent{ID: "a"}}},
		{"UpdateParty", TypeUpdateParty, &UpdatePartyBody{Context: TAPContext, Type: TypeUpdateParty, Party: &Party{ID: "p"}, Role: "originator"}},
		{"AddAgents", TypeAddAgents, &AddAgentsBody{Context: TAPContext, Type: TypeAddAgents, Agents: []Agent{{ID: "a"}}}},
		{"ReplaceAgent", TypeReplaceAgent, &ReplaceAgentBody{Context: TAPContext, Type: TypeReplaceAgent, Original: "old", Replacement: &Agent{ID: "new"}}},
		{"RemoveAgent", TypeRemoveAgent, &RemoveAgentBody{Context: TAPContext, Type: TypeRemoveAgent, Agent: "did:web:old"}},
		{"ConfirmRelationship", TypeConfirmRelationship, &ConfirmRelationshipBody{Context: TAPContext, Type: TypeAgent, ID: "did:pkh:eip155:1:0x1234", For: NewForField("did:web:beneficiary.vasp"), Role: "SettlementAddress"}},
		{"UpdatePolicies", TypeUpdatePolicies, &UpdatePoliciesBody{Context: TAPContext, Type: TypeUpdatePolicies, Policies: []Policy{{Type: "RequireAuthorization"}}}},
		{"Connect", TypeConnect, &ConnectBody{Context: TAPContext, Type: TypeConnect, Requester: &Party{ID: "r"}, Principal: &Party{ID: "p"}, Agents: []Agent{{ID: "a"}}, Constraints: &TransactionConstraints{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawBody, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			msg := &didcomm.Message{
				ID:   "test-id",
				Type: tt.msgType,
				Body: rawBody,
			}

			parsed, err := ParseBody(msg)
			if err != nil {
				t.Fatalf("ParseBody: %v", err)
			}
			if parsed.TAPType() != tt.msgType {
				t.Errorf("TAPType: got %q, want %q", parsed.TAPType(), tt.msgType)
			}
		})
	}
}

func TestParseBody_UnknownType(t *testing.T) {
	msg := &didcomm.Message{
		ID:   "test",
		Type: "https://example.com/unknown",
		Body: json.RawMessage(`{}`),
	}

	_, err := ParseBody(msg)
	if !errors.Is(err, ErrUnknownMessageType) {
		t.Errorf("expected ErrUnknownMessageType, got %v", err)
	}
}

func TestParseBody_InvalidJSON(t *testing.T) {
	msg := &didcomm.Message{
		ID:   "test",
		Type: TypeTransfer,
		Body: json.RawMessage(`{invalid`),
	}

	_, err := ParseBody(msg)
	if !errors.Is(err, ErrInvalidBody) {
		t.Errorf("expected ErrInvalidBody, got %v", err)
	}
}

func TestAllTypes_Count(t *testing.T) {
	types := AllTypes()
	if len(types) != 20 {
		t.Errorf("AllTypes: got %d types, want 20", len(types))
	}

	// Verify all are unique
	seen := make(map[string]bool)
	for _, typ := range types {
		if seen[typ] {
			t.Errorf("duplicate type: %s", typ)
		}
		seen[typ] = true
	}
}
