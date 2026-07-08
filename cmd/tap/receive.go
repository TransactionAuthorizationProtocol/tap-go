package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	tap "github.com/TransactionAuthorizationProtocol/tap-go"
)

// receiveOutput is the JSON output format for the receive command.
type receiveOutput struct {
	Message   json.RawMessage `json:"message"`
	Body      json.RawMessage `json:"body"`
	BodyType  string          `json:"bodyType"`
	SenderDID string          `json:"senderDid"`
	Encrypted bool            `json:"encrypted"`
	Anonymous bool            `json:"anonymous"`
}

func runReceive(args []string) error {
	fs := flag.NewFlagSet("receive", flag.ContinueOnError)
	keyFile := fs.String("key-file", "", "path to keys.json (required)")
	didDoc := fs.String("did-doc", "", "comma-separated DID document file paths")
	message := fs.String("message", "-", "message input: - (stdin), @file, or inline JSON")
	allowUnverified := fs.Bool("allow-unverified", false, "accept unauthenticated messages (sender NOT verified)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *keyFile == "" {
		return errors.New("--key-file is required")
	}

	dcClient, err := buildClient(*keyFile, *didDoc)
	if err != nil {
		return err
	}
	tapClient := tap.NewClient(dcClient)

	data, err := readInput(*message)
	if err != nil {
		return fmt.Errorf("read message: %w", err)
	}

	ctx := context.Background()
	receive := tapClient.Receive
	if *allowUnverified {
		receive = tapClient.ReceiveUnverified
	}
	result, err := receive(ctx, data)
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
		SenderDID: result.SenderDID,
		Encrypted: result.Encrypted,
		Anonymous: result.Anonymous,
	}
	outBytes, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}
	_, err = fmt.Fprintln(os.Stdout, string(outBytes))
	return err
}
