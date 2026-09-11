package tap

import (
	"errors"
	"testing"
)

// TAIP-5 marks `for` REQUIRED on every agent, but gives no way to say that who
// owns an agent is not established yet — a state real flows pass through. This
// package enforces `for` where the sender cannot honestly be unsure: the
// blockchain-address roles, whose owner whoever supplied the address knows.
//
// Inventing an owner is worse than omitting one: a receiver stores it as fact,
// and it is the difference between a self-hosted wallet (must prove ownership)
// and a custodied one (must not).

func TestAgent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		agent   Agent
		wantErr bool
	}{
		{"valid single for", Agent{ID: "did:web:vasp.example", For: NewForField("did:eg:bob")}, false},
		{"valid multiple for", Agent{ID: "did:web:shared.wallet", For: NewForField("did:web:vasp1", "did:web:vasp2")}, false},
		{"missing @id", Agent{For: NewForField("did:eg:bob")}, true},
		{"for with empty string", Agent{ID: "did:web:vasp.example", For: NewForField("")}, true},

		// An address is supplied by somebody who knows whose it is.
		{"settlement address without for", Agent{ID: "did:pkh:eip155:1:0xabc", Role: "SettlementAddress"}, true},
		{"source address without for", Agent{ID: "did:pkh:eip155:1:0xabc", Role: "SourceAddress"}, true},
		{"role casing is ignored", Agent{ID: "did:pkh:eip155:1:0xabc", Role: "settlementaddress"}, true},
		{
			name:  "settlement address with for",
			agent: Agent{ID: "did:pkh:eip155:1:0xabc", Role: "SettlementAddress", For: NewForField("did:web:vasp")},
		},

		// An institution may legitimately appear before anyone has said whose
		// behalf it acts on — an agent added mid-transfer is the common case.
		{"VASP without for", Agent{ID: "did:web:vasp.example", Role: "VASP"}, false},
		{"no role and no for", Agent{ID: "did:web:unknown.example"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.agent.Validate()
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.wantErr && !errors.Is(err, ErrInvalidBody) {
				t.Errorf("error should wrap ErrInvalidBody, got %v", err)
			}
		})
	}
}
