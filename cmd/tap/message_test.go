package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	tap "github.com/TransactionAuthorizationProtocol/tap-go"
	didcomm "github.com/notabene-id/go-didcomm"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "tap")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %s\n%s", err, out)
	}
	return binary
}

func TestCLI_Help(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "help")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("help failed: %s", err)
	}
	if !strings.Contains(string(out), "tap") {
		t.Fatal("help output missing 'tap'")
	}
	if !strings.Contains(string(out), "message") {
		t.Fatal("help output missing 'message'")
	}
	if !strings.Contains(string(out), "receive") {
		t.Fatal("help output missing 'receive'")
	}
}

func TestCLI_Version(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "version")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("version failed: %s", err)
	}
	if !strings.Contains(string(out), "tap") {
		t.Fatal("version output missing 'tap'")
	}
}

func TestCLI_UnknownCommand(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "foobar")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
	if !strings.Contains(string(out), "unknown command") {
		t.Fatalf("expected 'unknown command' error, got: %s", out)
	}
}

func TestCLI_NoArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestCLI_MessageNoType(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "message")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error with no message type")
	}
	if !strings.Contains(string(out), "message type required") {
		t.Fatalf("expected 'message type required' error, got: %s", out)
	}
}

func TestCLI_MessageUnknownType(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "message", "foobar", "--from", "did:key:z1", "--to", "did:key:z2")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for unknown message type")
	}
	if !strings.Contains(string(out), "unknown message type") {
		t.Fatalf("expected 'unknown message type' error, got: %s", out)
	}
}

func TestCLI_MessageMissingFrom(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "message", "transfer", "--to", "did:key:z2")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for missing --from")
	}
	if !strings.Contains(string(out), "--from is required") {
		t.Fatalf("expected '--from is required' error, got: %s", out)
	}
}

func TestCLI_MessageMissingTo(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "message", "transfer", "--from", "did:key:z1")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for missing --to")
	}
	if !strings.Contains(string(out), "--to is required") {
		t.Fatalf("expected '--to is required' error, got: %s", out)
	}
}

func TestCLI_MessageMissingThid(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "message", "authorize", "--from", "did:key:z1", "--to", "did:key:z2")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for missing --thid")
	}
	if !strings.Contains(string(out), "--thid is required") {
		t.Fatalf("expected '--thid is required' error, got: %s", out)
	}
}

func TestCLI_MessageTransfer(t *testing.T) {
	bin := buildBinary(t)
	body := `{"asset":"eip155:1/slip44:60","amount":"1.5","agents":[{"@id":"did:key:z1","role":"OriginatingVASP"}]}`

	cmd := exec.Command(bin, "message", "transfer",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		stderr := ""
		if ok {
			stderr = string(ee.Stderr)
		}
		t.Fatalf("message transfer failed: %s\n%s", err, stderr)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON output: %s\noutput: %s", err, out)
	}

	if msg.Type != tap.TypeTransfer {
		t.Fatalf("expected type %s, got %s", tap.TypeTransfer, msg.Type)
	}
	if msg.From != "did:key:z1" {
		t.Fatalf("expected from did:key:z1, got %s", msg.From)
	}
	if len(msg.To) != 1 || msg.To[0] != "did:key:z2" {
		t.Fatalf("expected to [did:key:z2], got %v", msg.To)
	}
	if msg.ID == "" {
		t.Fatal("expected non-empty message ID")
	}

	// Verify body has @context and @type
	var bodyMap map[string]any
	if err := json.Unmarshal(msg.Body, &bodyMap); err != nil {
		t.Fatal(err)
	}
	if bodyMap["@context"] != tap.TAPContext {
		t.Fatalf("expected @context %s, got %v", tap.TAPContext, bodyMap["@context"])
	}
	if bodyMap["@type"] != tap.TypeTransfer {
		t.Fatalf("expected @type %s, got %v", tap.TypeTransfer, bodyMap["@type"])
	}
}

func TestCLI_MessageTransfer_FileBody(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	bodyFile := filepath.Join(dir, "body.json")
	body := `{"asset":"eip155:1/slip44:60","amount":"2.0","agents":[{"@id":"did:key:z1","role":"OriginatingVASP"}]}`
	if err := os.WriteFile(bodyFile, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin, "message", "transfer",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--body", "@"+bodyFile,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message transfer with file body failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON output: %s", err)
	}
	if msg.Type != tap.TypeTransfer {
		t.Fatalf("expected type %s, got %s", tap.TypeTransfer, msg.Type)
	}
}

func TestCLI_MessagePayment(t *testing.T) {
	bin := buildBinary(t)
	body := `{"amount":"100","currency":"USD","merchant":{"@id":"did:key:z2"},"agents":[{"@id":"did:key:z1","role":"OriginatingVASP"}]}`

	cmd := exec.Command(bin, "message", "payment",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message payment failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypePayment {
		t.Fatalf("expected type %s, got %s", tap.TypePayment, msg.Type)
	}
}

func TestCLI_MessageRFQ(t *testing.T) {
	bin := buildBinary(t)
	body := `{"fromAssets":["eip155:1/slip44:60"],"toAssets":["eip155:1/slip44:0"],"fromAmount":"1.0","requester":{"@id":"did:key:z1"},"agents":[{"@id":"did:key:z1","role":"OriginatingVASP"}]}`

	cmd := exec.Command(bin, "message", "rfq",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message rfq failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeRFQ {
		t.Fatalf("expected type %s, got %s", tap.TypeRFQ, msg.Type)
	}
}

func TestCLI_MessageLock(t *testing.T) {
	bin := buildBinary(t)
	body := `{"asset":"eip155:1/slip44:60","amount":"5.0","originator":{"@id":"did:key:z1"},"beneficiary":{"@id":"did:key:z2"},"expiry":"2025-12-31T23:59:59Z","agents":[{"@id":"did:key:z1","role":"OriginatingVASP"}]}`

	cmd := exec.Command(bin, "message", "lock",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message lock failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeLock {
		t.Fatalf("expected type %s, got %s", tap.TypeLock, msg.Type)
	}
}

func TestCLI_MessageConnect(t *testing.T) {
	bin := buildBinary(t)
	body := `{"requester":{"@id":"did:key:z1"},"principal":{"@id":"did:key:z1"},"agents":[{"@id":"did:key:z1","role":"OriginatingVASP"}],"constraints":{"purposes":["travel"]}}`

	cmd := exec.Command(bin, "message", "connect",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message connect failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeConnect {
		t.Fatalf("expected type %s, got %s", tap.TypeConnect, msg.Type)
	}
}

func TestCLI_MessageAuthorize(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "message", "authorize",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message authorize failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeAuthorize {
		t.Fatalf("expected type %s, got %s", tap.TypeAuthorize, msg.Type)
	}
	if msg.Thid != "thread-123" {
		t.Fatalf("expected thid thread-123, got %s", msg.Thid)
	}
}

func TestCLI_MessageAuthorizationRequired(t *testing.T) {
	bin := buildBinary(t)
	body := `{"authorizationUrl":"https://example.com/auth","expires":"2025-12-31T23:59:59Z"}`

	cmd := exec.Command(bin, "message", "authorization-required",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message authorization-required failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeAuthorizationRequired {
		t.Fatalf("expected type %s, got %s", tap.TypeAuthorizationRequired, msg.Type)
	}
}

func TestCLI_MessageSettle(t *testing.T) {
	bin := buildBinary(t)
	body := `{"settlementAddress":"0x1234567890abcdef"}`

	cmd := exec.Command(bin, "message", "settle",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message settle failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeSettle {
		t.Fatalf("expected type %s, got %s", tap.TypeSettle, msg.Type)
	}
}

func TestCLI_MessageReject(t *testing.T) {
	bin := buildBinary(t)
	body := `{"reason":"compliance"}`

	cmd := exec.Command(bin, "message", "reject",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message reject failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeReject {
		t.Fatalf("expected type %s, got %s", tap.TypeReject, msg.Type)
	}
}

func TestCLI_MessageCancel(t *testing.T) {
	bin := buildBinary(t)
	body := `{"by":"did:key:z1"}`

	cmd := exec.Command(bin, "message", "cancel",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message cancel failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeCancel {
		t.Fatalf("expected type %s, got %s", tap.TypeCancel, msg.Type)
	}
}

func TestCLI_MessageRevert(t *testing.T) {
	bin := buildBinary(t)
	body := `{"settlementAddress":"0x1234","reason":"fraud"}`

	cmd := exec.Command(bin, "message", "revert",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message revert failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeRevert {
		t.Fatalf("expected type %s, got %s", tap.TypeRevert, msg.Type)
	}
}

func TestCLI_MessageCapture(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "message", "capture",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message capture failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeCapture {
		t.Fatalf("expected type %s, got %s", tap.TypeCapture, msg.Type)
	}
}

func TestCLI_MessageQuote(t *testing.T) {
	bin := buildBinary(t)
	body := `{"fromAsset":"eip155:1/slip44:60","toAsset":"eip155:1/slip44:0","fromAmount":"1.0","toAmount":"2000.0","provider":{"@id":"did:key:z1"},"agents":[{"@id":"did:key:z1","role":"Provider"}],"expiresAt":"2025-12-31T23:59:59Z"}`

	cmd := exec.Command(bin, "message", "quote",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message quote failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeQuote {
		t.Fatalf("expected type %s, got %s", tap.TypeQuote, msg.Type)
	}
}

func TestCLI_MessageAddAgents(t *testing.T) {
	bin := buildBinary(t)
	body := `{"agents":[{"@id":"did:key:z3","role":"IntermediaryVASP"}]}`

	cmd := exec.Command(bin, "message", "add-agents",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message add-agents failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeAddAgents {
		t.Fatalf("expected type %s, got %s", tap.TypeAddAgents, msg.Type)
	}
}

func TestCLI_MessageRemoveAgent(t *testing.T) {
	bin := buildBinary(t)
	body := `{"agent":"did:key:z3"}`

	cmd := exec.Command(bin, "message", "remove-agent",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message remove-agent failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeRemoveAgent {
		t.Fatalf("expected type %s, got %s", tap.TypeRemoveAgent, msg.Type)
	}
}

func TestCLI_MessageReplaceAgent(t *testing.T) {
	bin := buildBinary(t)
	body := `{"original":"did:key:z3","replacement":{"@id":"did:key:z4","role":"IntermediaryVASP"}}`

	cmd := exec.Command(bin, "message", "replace-agent",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message replace-agent failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeReplaceAgent {
		t.Fatalf("expected type %s, got %s", tap.TypeReplaceAgent, msg.Type)
	}
}

func TestCLI_MessageUpdateAgent(t *testing.T) {
	bin := buildBinary(t)
	body := `{"agent":{"@id":"did:key:z3","role":"IntermediaryVASP"}}`

	cmd := exec.Command(bin, "message", "update-agent",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message update-agent failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeUpdateAgent {
		t.Fatalf("expected type %s, got %s", tap.TypeUpdateAgent, msg.Type)
	}
}

func TestCLI_MessageUpdateParty(t *testing.T) {
	bin := buildBinary(t)
	body := `{"party":{"@id":"did:key:z3","name":"Alice"},"role":"originator"}`

	cmd := exec.Command(bin, "message", "update-party",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message update-party failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeUpdateParty {
		t.Fatalf("expected type %s, got %s", tap.TypeUpdateParty, msg.Type)
	}
}

func TestCLI_MessageUpdatePolicies(t *testing.T) {
	bin := buildBinary(t)
	body := `{"policies":[{"@type":"RequirePresentation","purpose":"compliance"}]}`

	cmd := exec.Command(bin, "message", "update-policies",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message update-policies failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeUpdatePolicies {
		t.Fatalf("expected type %s, got %s", tap.TypeUpdatePolicies, msg.Type)
	}
}

func TestCLI_MessageConfirmRelationship(t *testing.T) {
	bin := buildBinary(t)
	body := `{"@id":"did:pkh:eip155:1:0x1234","for":"did:key:z1","role":"SettlementAddress"}`

	cmd := exec.Command(bin, "message", "confirm-relationship",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--thid", "thread-123",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message confirm-relationship failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if msg.Type != tap.TypeConfirmRelationship {
		t.Fatalf("expected type %s, got %s", tap.TypeConfirmRelationship, msg.Type)
	}
}

// TestCLI_MessageTransfer_MultipleRecipients tests comma-separated --to flag.
func TestCLI_MessageTransfer_MultipleRecipients(t *testing.T) {
	bin := buildBinary(t)
	body := `{"asset":"eip155:1/slip44:60","amount":"1.0","agents":[{"@id":"did:key:z1","role":"OriginatingVASP"}]}`

	cmd := exec.Command(bin, "message", "transfer",
		"--from", "did:key:z1",
		"--to", "did:key:z2,did:key:z3",
		"--body", body,
	)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("message transfer with multiple recipients failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(out, &msg); err != nil {
		t.Fatalf("invalid JSON: %s", err)
	}
	if len(msg.To) != 2 {
		t.Fatalf("expected 2 recipients, got %d", len(msg.To))
	}
	if msg.To[0] != "did:key:z2" || msg.To[1] != "did:key:z3" {
		t.Fatalf("unexpected recipients: %v", msg.To)
	}
}

// TestCLI_DIDGenerateKey tests that DID commands work through the tap binary.
// DID generation and envelope pack/unpack moved to the didcomm CLI, so the tap
// CLI no longer exposes `did`/`pack`/`unpack`/`send`.

// TestCLI_MessageTransfer_BodyValidation tests that body validation errors propagate.
func TestCLI_MessageTransfer_BodyValidation(t *testing.T) {
	bin := buildBinary(t)
	// Missing required "asset" field
	body := `{"amount":"1.0","agents":[{"@id":"did:key:z1"}]}`

	cmd := exec.Command(bin, "message", "transfer",
		"--from", "did:key:z1",
		"--to", "did:key:z2",
		"--body", body,
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected validation error for missing asset")
	}
	if !strings.Contains(string(out), "asset") {
		t.Fatalf("expected asset validation error, got: %s", out)
	}
}
