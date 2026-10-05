package tap

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	didcomm "github.com/notabene-id/go-didcomm"
)

func TestNewRejectMessage(t *testing.T) {
	body := &RejectBody{Reason: "Beneficiary name mismatch"}
	msg, err := NewRejectMessage("did:web:beneficiary.vasp", []string{"did:web:originator.vasp"}, "1234567890", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Type != TypeReject {
		t.Errorf("Type: got %q", msg.Type)
	}
}

func TestNewRejectMessage_NoReasonOptional(t *testing.T) {
	body := &RejectBody{}
	_, err := NewRejectMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRejectBody_JSONRoundTrip(t *testing.T) {
	body := RejectBody{
		Context: TAPContext,
		Type:    TypeReject,
		Reason:  "Compliance check failed",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got RejectBody
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Reason != body.Reason {
		t.Errorf("mismatch: %+v", got)
	}
}

func TestReject_TestVectorValid(t *testing.T) {
	data, err := os.ReadFile("TAIPs/test-vectors/reject/valid.json")
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

	var body RejectBody
	if err := json.Unmarshal(tv.Message.Body, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if body.Reason != "Beneficiary name mismatch" {
		t.Errorf("Reason: got %q", body.Reason)
	}
}

func TestRejectBody_ParseBody(t *testing.T) {
	body := &RejectBody{Reason: "test"}
	msg, err := NewRejectMessage("from", []string{"to"}, "thid", body)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	parsed, err := ParseBody(msg)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rb, ok := parsed.(*RejectBody)
	if !ok {
		t.Fatalf("expected *RejectBody, got %T", parsed)
	}
	if rb.Reason != "test" {
		t.Errorf("Reason: got %q", rb.Reason)
	}
}

// Code carries the ISO 20022 status reason next to the free-text reason, and is
// left out entirely when unset so a Reject without one reads exactly as before.
func TestRejectBody_Code(t *testing.T) {
	cases := []struct {
		name     string
		body     RejectBody
		wantJSON string
	}{
		{
			name:     "code_round_trips",
			body:     RejectBody{Reason: "Blockchain address is not ours", Code: RejectCodeInvalidCreditorAccountNumber},
			wantJSON: `"code":"AC03"`,
		},
		{
			name: "empty_code_is_omitted",
			body: RejectBody{Reason: "Compliance check failed"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := NewRejectMessage("from", []string{"to"}, "thid", &tc.body)
			if err != nil {
				t.Fatalf("create: %v", err)
			}
			raw := string(msg.Body)
			if tc.wantJSON == "" && strings.Contains(raw, `"code"`) {
				t.Fatalf("empty code must be omitted, got %s", raw)
			}
			if tc.wantJSON != "" && !strings.Contains(raw, tc.wantJSON) {
				t.Fatalf("body %s does not contain %s", raw, tc.wantJSON)
			}
			parsed, err := ParseBody(msg)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			rb, ok := parsed.(*RejectBody)
			if !ok {
				t.Fatalf("expected *RejectBody, got %T", parsed)
			}
			if rb.Code != tc.body.Code {
				t.Errorf("Code: got %q, want %q", rb.Code, tc.body.Code)
			}
		})
	}
}

// Sending is strict: only a code this library defines goes out, so a typo or an
// invented code fails at construction. Reading is lenient: an unknown code from
// another implementation is kept for the receiver to ignore, not a parse error.
func TestRejectCode_SendStrictReadLenient(t *testing.T) {
	cases := []struct {
		name    string
		code    RejectCode
		wantErr bool
	}{
		{name: "known_code_is_sent", code: RejectCodeInvalidCreditorAccountNumber},
		{name: "name_mismatch_code_is_sent", code: RejectCodeInconsistentWithEndCustomer},
		{name: "no_code_is_sent", code: ""},
		{name: "unknown_code_is_refused", code: "AC99", wantErr: true},
		{name: "lowercase_known_code_is_refused", code: "ac03", wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewRejectMessage("from", []string{"to"}, "thid", &RejectBody{Reason: "r", Code: tc.code})
			if tc.wantErr {
				if !errors.Is(err, ErrInvalidBody) {
					t.Fatalf("want ErrInvalidBody, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}

	t.Run("unknown_code_is_kept_on_read", func(t *testing.T) {
		msg := &didcomm.Message{
			Type: TypeReject,
			Body: json.RawMessage(`{"@context":"https://tap.rsvp/schema/1.0",` +
				`"@type":"https://tap.rsvp/schema/1.0#Reject","reason":"r","code":"AC04"}`),
		}
		parsed, err := ParseBody(msg)
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		if got := parsed.(*RejectBody).Code; got != "AC04" {
			t.Errorf("Code: got %q, want AC04", got)
		}
	})
}
