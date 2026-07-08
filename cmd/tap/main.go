package main

import (
	"fmt"
	"os"
)

const version = "0.5.0"

const usage = `tap - TAP (Transaction Authorization Protocol) CLI

Usage:
  tap <command> [options]

Commands:
  message <type> --from <did> --to <did> [flags]            Create a TAP message
  receive        --key-file <f> [--did-doc <f>]             Unpack + parse TAP body
  version                                                   Print version
  help                                                      Print this help

For DID generation and envelope pack/unpack/send use the "didcomm" CLI
(github.com/notabene-id/go-didcomm/cmd/didcomm).

TAP message types:
  Initiating:  transfer, payment, rfq, lock, connect
  Reply:       authorize, authorization-required, settle, reject, cancel,
               revert, capture, quote, add-agents, remove-agent,
               replace-agent, update-agent, update-party, update-policies,
               confirm-relationship

Message flags:
  --from <did>     Sender DID (required)
  --to <did>       Recipient DID(s), comma-separated (required)
  --thid <id>      Thread ID (required for reply messages)
  --body <json>    Body JSON: inline, @file.json, or - for stdin (default: {})

Message input (--message flag):
  -           Read from stdin (default)
  @file.json  Read from file
  '{"json"}'  Inline JSON string

Examples:
  # Generate identities
  tap did generate-key --output-dir alice
  tap did generate-key --output-dir bob

  # Create a TAP transfer message
  ALICE=$(jq -r .id alice/did-doc.json)
  BOB=$(jq -r .id bob/did-doc.json)
  tap message transfer --from $ALICE --to $BOB \
    --body '{"asset":"eip155:1/slip44:60","amount":"1.5","agents":[{"@id":"'$ALICE'","role":"OriginatingVASP"}]}'

  # Pipe: create message → pack (via the didcomm CLI) → send
  tap message transfer --from $ALICE --to $BOB --body @body.json | \
    didcomm pack --keys alice/keys.json --profile 1pu-v3

  # Receive: unpack + parse TAP body
  echo '<packed-message>' | tap receive --key-file bob/keys.json
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "message":
		err = runMessage(os.Args[2:])
	case "receive":
		err = runReceive(os.Args[2:])
	case "version":
		fmt.Println("tap " + version)
	case "help", "--help", "-h":
		fmt.Print(usage)
	default:
		fmt.Fprintln(os.Stderr, "unknown command: "+os.Args[1]+"\n") //nolint:gosec // CLI stderr output
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
