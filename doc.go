// Package tap provides typed Go wrappers for all 20 TAP (Transaction Authorization
// Protocol) message types, built on top of the go-didcomm library.
//
// TAP defines a DIDComm v2-based protocol for multi-party transaction authorization
// before blockchain settlement. This package provides:
//
//   - Typed body structs for each message type (TransferBody, PaymentBody, etc.)
//   - Constructor functions that validate required fields and produce didcomm.Message values
//   - A ParseBody dispatcher that unmarshals any TAP message body into its typed struct
//   - A Client wrapper around didcomm.Client that combines unpacking and parsing
//
// # Message Types
//
// TAP defines 20 message types organized into four categories:
//
// Transaction messages initiate financial operations:
//   - Transfer (TAIP-3) — asset transfer between parties
//   - Payment (TAIP-14) — merchant payment request
//   - Exchange (TAIP-18) — asset exchange request
//   - Quote (TAIP-18) — exchange price quote response
//   - Escrow (TAIP-17) — hold funds in escrow
//
// Authorization flow messages manage transaction lifecycle:
//   - Authorize (TAIP-4) — approve a transaction
//   - AuthorizationRequired (TAIP-4) — request interactive authorization
//   - Settle (TAIP-4) — announce blockchain settlement
//   - Reject (TAIP-4) — reject a transaction
//   - Cancel (TAIP-4) — cancel a transaction
//   - Revert (TAIP-4) — request reversal of settled transaction
//   - Capture (TAIP-17) — release escrowed funds
//
// Agent management messages modify transaction participants:
//   - UpdateAgent (TAIP-5) — update agent information
//   - UpdateParty (TAIP-6) — update party information
//   - AddAgents (TAIP-5) — add new agents
//   - ReplaceAgent (TAIP-5) — replace an existing agent
//   - RemoveAgent (TAIP-5) — remove an agent
//
// Relationship messages manage agent connections:
//   - ConfirmRelationship (TAIP-9) — confirm a relationship between parties
//   - UpdatePolicies (TAIP-7) — update agent policies
//   - Connect (TAIP-15) — establish agent connection with constraints
//
// # Usage
//
// Create a message with a constructor:
//
//	msg, err := tap.NewTransferMessage(
//	    "did:web:originator.vasp",
//	    []string{"did:web:beneficiary.vasp"},
//	    &tap.TransferBody{
//	        Asset:  "eip155:1/slip44:60",
//	        Amount: "1.23",
//	        Agents: []tap.Agent{{ID: "did:web:originator.vasp", For: tap.NewForField("did:eg:bob")}},
//	    },
//	)
//
// Parse a received message:
//
//	body, err := tap.ParseBody(msg)
//	switch b := body.(type) {
//	case *tap.TransferBody:
//	    fmt.Printf("Transfer of %s %s\n", b.Amount, b.Asset)
//	case *tap.PaymentBody:
//	    fmt.Printf("Payment of %s\n", b.Amount)
//	}
//
// Use the Client for end-to-end DIDComm + TAP parsing:
//
//	client := tap.NewClient(didcommClient)
//	result, err := client.Receive(ctx, envelope)
//	// result.Body is a typed TAPBody, result.Encrypted/Signed are set
//
// # Specifications
//
// TAP is defined by the TAIPs (Transaction Authorization Improvement Proposals):
// https://tap.rsvp
//
// All message bodies use the JSON-LD context "https://tap.rsvp/schema/1.0" and
// type URLs of the form "https://tap.rsvp/schema/1.0#MessageType".
package tap
