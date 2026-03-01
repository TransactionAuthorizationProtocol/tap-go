package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Notabene-id/go-didcomm/cli"

	tap "github.com/TransactionAuthorizationProtocol/tap-go"
)

// receiveOutput is the JSON output format for the receive command.
type receiveOutput struct {
	Message   json.RawMessage `json:"message"`
	Body      json.RawMessage `json:"body"`
	BodyType  string          `json:"bodyType"`
	Encrypted bool            `json:"encrypted"`
	Signed    bool            `json:"signed"`
	Anonymous bool            `json:"anonymous"`
}

func runReceive(args []string) error {
	fs := flag.NewFlagSet("receive", flag.ContinueOnError)
	keyFile := fs.String("key-file", "", "path to JWK Set file with private keys (required)")
	didDoc := fs.String("did-doc", "", "comma-separated DID document file paths")
	message := fs.String("message", "-", "message input: - (stdin), @file, or inline JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *keyFile == "" {
		return fmt.Errorf("--key-file is required")
	}

	dcClient, err := cli.BuildClient(*keyFile, *didDoc)
	if err != nil {
		return err
	}

	tapClient := tap.NewClient(dcClient)

	data, err := cli.ReadMessageInput(*message)
	if err != nil {
		return fmt.Errorf("read message: %w", err)
	}

	result, err := tapClient.Receive(context.Background(), data)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}

	msgBytes, err := json.Marshal(result.Message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	bodyBytes, err := json.Marshal(result.Body)
	if err != nil {
		return fmt.Errorf("marshal body: %w", err)
	}

	out := receiveOutput{
		Message:   msgBytes,
		Body:      bodyBytes,
		BodyType:  result.Body.TAPType(),
		Encrypted: result.Encrypted,
		Signed:    result.Signed,
		Anonymous: result.Anonymous,
	}

	outBytes, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}

	_, err = fmt.Fprintln(os.Stdout, string(outBytes))
	return err
}
