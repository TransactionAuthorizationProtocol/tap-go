package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	didcomm "github.com/Notabene-id/go-didcomm"
	"github.com/Notabene-id/go-didcomm/cli"
	tap "github.com/TransactionAuthorizationProtocol/tap-go"
)

// generateIdentity creates a DID identity and writes files to dir.
func generateIdentity(t *testing.T, dir string) *didcomm.DIDDocument {
	t.Helper()
	doc, kp, err := didcomm.GenerateDIDKey()
	if err != nil {
		t.Fatal(err)
	}

	docBytes, err := cli.MarshalDIDDoc(doc)
	if err != nil {
		t.Fatal(err)
	}
	keyBytes, err := cli.MarshalKeyPair(kp)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "did-doc.json"), docBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "keys.json"), keyBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestCLI_ReceiveMissingKeyFile(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "receive", "--message", "{}")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for missing key file")
	}
	if !strings.Contains(string(out), "--key-file is required") {
		t.Fatalf("expected key-file error, got: %s", out)
	}
}

func TestCLI_Receive_SignedTransfer(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	aliceDir := filepath.Join(dir, "alice")
	aliceDoc := generateIdentity(t, aliceDir)

	// Create a TAP transfer message
	transferBody := &tap.TransferBody{
		Asset:  "eip155:1/slip44:60",
		Amount: "1.5",
		Agents: []tap.Agent{{ID: aliceDoc.ID, Role: "OriginatingVASP"}},
	}
	msg, err := tap.NewTransferMessage(aliceDoc.ID, []string{aliceDoc.ID}, transferBody)
	if err != nil {
		t.Fatal(err)
	}

	// Pack it signed
	client, err := cli.BuildClient(
		filepath.Join(aliceDir, "keys.json"),
		filepath.Join(aliceDir, "did-doc.json"),
	)
	if err != nil {
		t.Fatal(err)
	}

	packed, err := client.PackSigned(context.Background(), msg)
	if err != nil {
		t.Fatal(err)
	}

	// Receive via CLI
	receiveCmd := exec.Command(bin, "receive",
		"--key-file", filepath.Join(aliceDir, "keys.json"),
		"--did-doc", filepath.Join(aliceDir, "did-doc.json"),
		"--message", string(packed),
	)
	out, err := receiveCmd.Output()
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		stderr := ""
		if ok {
			stderr = string(ee.Stderr)
		}
		t.Fatalf("receive failed: %s\n%s", err, stderr)
	}

	var result receiveOutput
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("invalid receive output: %s\noutput: %s", err, out)
	}

	if !result.Signed {
		t.Fatal("expected signed=true")
	}
	if result.Encrypted {
		t.Fatal("expected encrypted=false")
	}
	if result.BodyType != tap.TypeTransfer {
		t.Fatalf("expected bodyType %s, got %s", tap.TypeTransfer, result.BodyType)
	}

	// Verify body content
	var body tap.TransferBody
	if err := json.Unmarshal(result.Body, &body); err != nil {
		t.Fatal(err)
	}
	if body.Asset != "eip155:1/slip44:60" {
		t.Fatalf("expected asset eip155:1/slip44:60, got %s", body.Asset)
	}
	if body.Amount != "1.5" {
		t.Fatalf("expected amount 1.5, got %s", body.Amount)
	}
}

func TestCLI_Receive_AuthcryptReject(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	aliceDir := filepath.Join(dir, "alice")
	aliceDoc := generateIdentity(t, aliceDir)
	bobDir := filepath.Join(dir, "bob")
	bobDoc := generateIdentity(t, bobDir)

	didDocs := filepath.Join(aliceDir, "did-doc.json") + "," + filepath.Join(bobDir, "did-doc.json")

	// Create a TAP reject message
	rejectBody := &tap.RejectBody{
		Reason: "compliance failure",
	}
	msg, err := tap.NewRejectMessage(aliceDoc.ID, []string{bobDoc.ID}, "thread-456", rejectBody)
	if err != nil {
		t.Fatal(err)
	}

	// Pack authcrypt
	client, err := cli.BuildClient(filepath.Join(aliceDir, "keys.json"), didDocs)
	if err != nil {
		t.Fatal(err)
	}

	packed, err := client.PackAuthcrypt(context.Background(), msg)
	if err != nil {
		t.Fatal(err)
	}

	// Receive via CLI with bob's keys
	receiveCmd := exec.Command(bin, "receive",
		"--key-file", filepath.Join(bobDir, "keys.json"),
		"--did-doc", didDocs,
		"--message", string(packed),
	)
	out, err := receiveCmd.Output()
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		stderr := ""
		if ok {
			stderr = string(ee.Stderr)
		}
		t.Fatalf("receive failed: %s\n%s", err, stderr)
	}

	var result receiveOutput
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("invalid receive output: %s", err)
	}

	if !result.Encrypted {
		t.Fatal("expected encrypted=true")
	}
	if !result.Signed {
		t.Fatal("expected signed=true")
	}
	if result.Anonymous {
		t.Fatal("expected anonymous=false")
	}
	if result.BodyType != tap.TypeReject {
		t.Fatalf("expected bodyType %s, got %s", tap.TypeReject, result.BodyType)
	}

	var body tap.RejectBody
	if err := json.Unmarshal(result.Body, &body); err != nil {
		t.Fatal(err)
	}
	if body.Reason != "compliance failure" {
		t.Fatalf("expected reason 'compliance failure', got %s", body.Reason)
	}
}

// TestCLI_MessagePipe_Transfer tests creating a message and piping to pack/unpack/receive.
func TestCLI_MessagePipe_Transfer(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()

	aliceDir := filepath.Join(dir, "alice")
	aliceDoc := generateIdentity(t, aliceDir)

	// Step 1: Create TAP message via CLI
	body := `{"asset":"eip155:1/slip44:60","amount":"3.0","agents":[{"@id":"` + aliceDoc.ID + `","role":"OriginatingVASP"}]}`
	msgCmd := exec.Command(bin, "message", "transfer",
		"--from", aliceDoc.ID,
		"--to", aliceDoc.ID,
		"--body", body,
	)
	msgOut, err := msgCmd.Output()
	if err != nil {
		t.Fatalf("message creation failed: %s", err)
	}

	// Step 2: Pack signed
	packCmd := exec.Command(bin, "pack", "signed",
		"--key-file", filepath.Join(aliceDir, "keys.json"),
		"--did-doc", filepath.Join(aliceDir, "did-doc.json"),
		"--message", string(msgOut),
	)
	packed, err := packCmd.Output()
	if err != nil {
		t.Fatalf("pack failed: %s", err)
	}

	// Step 3: Receive (unpack + TAP parse)
	receiveCmd := exec.Command(bin, "receive",
		"--key-file", filepath.Join(aliceDir, "keys.json"),
		"--did-doc", filepath.Join(aliceDir, "did-doc.json"),
		"--message", string(packed),
	)
	receiveOut, err := receiveCmd.Output()
	if err != nil {
		t.Fatalf("receive failed: %s", err)
	}

	var result receiveOutput
	if err := json.Unmarshal(receiveOut, &result); err != nil {
		t.Fatalf("invalid receive output: %s", err)
	}

	if result.BodyType != tap.TypeTransfer {
		t.Fatalf("expected bodyType %s, got %s", tap.TypeTransfer, result.BodyType)
	}
	if !result.Signed {
		t.Fatal("expected signed=true")
	}
}
