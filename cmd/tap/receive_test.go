package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	didcomm "github.com/notabene-id/go-didcomm"
	"github.com/notabene-id/go-didcomm/softkey"

	tap "github.com/TransactionAuthorizationProtocol/tap-go"
)

// generateIdentity creates a did:key identity, writes keys.json + did-doc.json
// to dir, and returns both.
func generateIdentity(t *testing.T, dir string) (*didcomm.DIDDocument, *didcomm.KeyMaterial) {
	t.Helper()
	doc, km, err := didcomm.GenerateDIDKey()
	if err != nil {
		t.Fatal(err)
	}
	docBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	keyBytes, err := json.MarshalIndent(km, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "did-doc.json"), docBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "keys.json"), keyBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	return doc, km
}

// packMessage packs msg in-process with the given profile, using a did:key
// auto-resolving client. did:key identities need no DID-document overrides.
func packMessage(t *testing.T, km *didcomm.KeyMaterial, msg *didcomm.Message, profile didcomm.Profile) []byte {
	t.Helper()
	store, err := softkey.New(km)
	if err != nil {
		t.Fatal(err)
	}
	resolver, _ := didcomm.DefaultResolver()
	packed, err := didcomm.NewClient(resolver, store).Pack(context.Background(), msg, didcomm.WithProfile(profile))
	if err != nil {
		t.Fatalf("pack: %v", err)
	}
	return packed
}

// receiveViaCLI runs `tap receive` for the given key file and packed message.
func receiveViaCLI(t *testing.T, bin, keyFile string, packed []byte) receiveOutput {
	t.Helper()
	cmd := exec.Command(bin, "receive", "--key-file", keyFile, "--message", string(packed))
	out, err := cmd.Output()
	if err != nil {
		stderr := ""
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = string(ee.Stderr)
		}
		t.Fatalf("receive failed: %s\n%s", err, stderr)
	}
	var result receiveOutput
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("invalid receive output: %s\noutput: %s", err, out)
	}
	return result
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
	aliceDir := filepath.Join(t.TempDir(), "alice")
	aliceDoc, aliceKM := generateIdentity(t, aliceDir)

	msg, err := tap.NewTransferMessage(aliceDoc.ID, []string{aliceDoc.ID}, &tap.TransferBody{
		Asset:  "eip155:1/slip44:60",
		Amount: "1.5",
		Agents: []tap.Agent{{ID: aliceDoc.ID, Role: "OriginatingVASP"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	packed := packMessage(t, aliceKM, msg, didcomm.ProfileSigned)

	result := receiveViaCLI(t, bin, filepath.Join(aliceDir, "keys.json"), packed)
	if result.SenderDID != aliceDoc.ID {
		t.Fatalf("senderDid = %q, want %q", result.SenderDID, aliceDoc.ID)
	}
	if result.Encrypted {
		t.Fatal("expected encrypted=false")
	}
	if result.BodyType != tap.TypeTransfer {
		t.Fatalf("bodyType = %s, want %s", result.BodyType, tap.TypeTransfer)
	}

	var body tap.TransferBody
	if err := json.Unmarshal(result.Body, &body); err != nil {
		t.Fatal(err)
	}
	if body.Asset != "eip155:1/slip44:60" || body.Amount != "1.5" {
		t.Fatalf("body = %+v", body)
	}
}

func TestCLI_Receive_AuthcryptReject(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	aliceDoc, aliceKM := generateIdentity(t, filepath.Join(dir, "alice"))
	bobDoc, _ := generateIdentity(t, filepath.Join(dir, "bob"))

	msg, err := tap.NewRejectMessage(aliceDoc.ID, []string{bobDoc.ID}, "thread-456", &tap.RejectBody{
		Reason: "compliance failure",
	})
	if err != nil {
		t.Fatal(err)
	}
	packed := packMessage(t, aliceKM, msg, didcomm.ProfileAuthcrypt1PUv3)

	result := receiveViaCLI(t, bin, filepath.Join(dir, "bob", "keys.json"), packed)
	if !result.Encrypted {
		t.Fatal("expected encrypted=true")
	}
	if result.SenderDID != aliceDoc.ID {
		t.Fatalf("senderDid = %q, want %q", result.SenderDID, aliceDoc.ID)
	}
	if result.Anonymous {
		t.Fatal("expected anonymous=false")
	}
	if result.BodyType != tap.TypeReject {
		t.Fatalf("bodyType = %s, want %s", result.BodyType, tap.TypeReject)
	}

	var body tap.RejectBody
	if err := json.Unmarshal(result.Body, &body); err != nil {
		t.Fatal(err)
	}
	if body.Reason != "compliance failure" {
		t.Fatalf("reason = %s", body.Reason)
	}
}

// TestCLI_MessagePipe_Transfer creates a message via the CLI, packs it
// in-process, and receives it via the CLI.
func TestCLI_MessagePipe_Transfer(t *testing.T) {
	bin := buildBinary(t)
	aliceDir := filepath.Join(t.TempDir(), "alice")
	aliceDoc, aliceKM := generateIdentity(t, aliceDir)

	body := `{"asset":"eip155:1/slip44:60","amount":"3.0","agents":[{"@id":"` + aliceDoc.ID + `","role":"OriginatingVASP"}]}`
	msgOut, err := exec.Command(bin, "message", "transfer",
		"--from", aliceDoc.ID, "--to", aliceDoc.ID, "--body", body,
	).Output()
	if err != nil {
		t.Fatalf("message creation failed: %s", err)
	}

	var msg didcomm.Message
	if err := json.Unmarshal(msgOut, &msg); err != nil {
		t.Fatalf("parse message: %s", err)
	}
	packed := packMessage(t, aliceKM, &msg, didcomm.ProfileSigned)

	result := receiveViaCLI(t, bin, filepath.Join(aliceDir, "keys.json"), packed)
	if result.BodyType != tap.TypeTransfer {
		t.Fatalf("bodyType = %s, want %s", result.BodyType, tap.TypeTransfer)
	}
	if result.SenderDID != aliceDoc.ID {
		t.Fatalf("senderDid = %q, want %q", result.SenderDID, aliceDoc.ID)
	}
}
