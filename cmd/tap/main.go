package main

import (
	"fmt"
	"os"

	"github.com/Notabene-id/go-didcomm/cli"
)

const version = "0.1.0"

const usage = `tap - TAP (Transaction Authorization Protocol) CLI

Usage:
  tap <command> [options]

Commands:
  did generate-key                                          Generate a did:key identity
  did generate-web --domain <d> [--path <p>]                Generate a did:web identity
  pack signed    --key-file <f> [--send] [--did-doc <f>]    Sign a message (JWS)
  pack anoncrypt [--send] [--did-doc <f>] [--message <m>]   Anonymous encrypt (JWE)
  pack authcrypt --key-file <f> [--send] [--did-doc <f>]    Sign-then-encrypt
  unpack         --key-file <f> [--did-doc <f>]             Unpack a message
  send           --to <url> [--message <m>]                 HTTP POST pre-packed message
  message <type> --from <did> --to <did> [flags]            Create a TAP message
  receive        --key-file <f> [--did-doc <f>]             Unpack + parse TAP body
  version                                                   Print version
  help                                                      Print this help

TAP message types:
  Initiating:  transfer, payment, exchange, escrow, connect
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

  # Pipe: create message → pack → send
  tap message transfer --from $ALICE --to $BOB --body @body.json | \
    tap pack authcrypt --key-file alice/keys.json

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
	case "did":
		err = cli.RunDID(os.Args[2:])
	case "pack":
		err = cli.RunPack(os.Args[2:])
	case "unpack":
		err = cli.RunUnpack(os.Args[2:])
	case "send":
		err = cli.RunSend(os.Args[2:])
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
