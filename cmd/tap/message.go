package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	didcomm "github.com/notabene-id/go-didcomm"

	tap "github.com/TransactionAuthorizationProtocol/tap-go"
)

const messageUsage = `Usage: tap message <type> --from <did> --to <did> [--thid <id>] [--body <json>]

Initiating types (no --thid):
  transfer, payment, rfq, lock, connect

Reply types (require --thid):
  authorize, authorization-required, settle, reject, cancel, revert,
  capture, quote, add-agents, remove-agent, replace-agent, update-agent,
  update-party, update-policies, confirm-relationship

Body input (--body flag):
  '{...}'     Inline JSON
  @file.json  Read from file
  -           Read from stdin
  (omitted)   Defaults to {}
`

// messageFlags holds the common flags for all message subcommands.
type messageFlags struct {
	from string
	to   []string
	thid string
	body []byte
}

func parseMessageFlags(name string, args []string, requireThid bool) (*messageFlags, error) {
	fs := flag.NewFlagSet("message "+name, flag.ContinueOnError)
	from := fs.String("from", "", "sender DID (required)")
	to := fs.String("to", "", "recipient DID(s), comma-separated (required)")
	thid := fs.String("thid", "", "thread ID (required for reply messages)")
	bodyFlag := fs.String("body", "", "body JSON: inline, @file.json, or - for stdin")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *from == "" {
		return nil, fmt.Errorf("--from is required")
	}
	if *to == "" {
		return nil, fmt.Errorf("--to is required")
	}
	if requireThid && *thid == "" {
		return nil, fmt.Errorf("--thid is required for %s messages", name)
	}

	// Parse recipient list
	var recipients []string
	for _, r := range strings.Split(*to, ",") {
		r = strings.TrimSpace(r)
		if r != "" {
			recipients = append(recipients, r)
		}
	}

	// Parse body
	var body []byte
	if *bodyFlag == "" {
		body = []byte("{}")
	} else {
		var err error
		body, err = readInput(*bodyFlag)
		if err != nil {
			return nil, fmt.Errorf("read body: %w", err)
		}
	}

	return &messageFlags{
		from: *from,
		to:   recipients,
		thid: *thid,
		body: body,
	}, nil
}

func writeMessage(msg *didcomm.Message) error {
	data, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	_, err = os.Stdout.Write(data)
	return err
}

func runMessage(args []string) error {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, messageUsage)
		return fmt.Errorf("message type required")
	}

	msgType := args[0]
	msgArgs := args[1:]

	switch msgType {
	// Initiating messages (no thid)
	case "transfer":
		return runMessageTransfer(msgArgs)
	case "payment":
		return runMessagePayment(msgArgs)
	case "rfq":
		return runMessageRFQ(msgArgs)
	case "lock":
		return runMessageLock(msgArgs)
	case "connect":
		return runMessageConnect(msgArgs)

	// Reply messages (require thid)
	case "authorize":
		return runMessageAuthorize(msgArgs)
	case "authorization-required":
		return runMessageAuthorizationRequired(msgArgs)
	case "settle":
		return runMessageSettle(msgArgs)
	case "reject":
		return runMessageReject(msgArgs)
	case "cancel":
		return runMessageCancel(msgArgs)
	case "revert":
		return runMessageRevert(msgArgs)
	case "capture":
		return runMessageCapture(msgArgs)
	case "quote":
		return runMessageQuote(msgArgs)
	case "add-agents":
		return runMessageAddAgents(msgArgs)
	case "remove-agent":
		return runMessageRemoveAgent(msgArgs)
	case "replace-agent":
		return runMessageReplaceAgent(msgArgs)
	case "update-agent":
		return runMessageUpdateAgent(msgArgs)
	case "update-party":
		return runMessageUpdateParty(msgArgs)
	case "update-policies":
		return runMessageUpdatePolicies(msgArgs)
	case "confirm-relationship":
		return runMessageConfirmRelationship(msgArgs)

	default:
		return fmt.Errorf("unknown message type: %s", msgType)
	}
}

// --- Initiating messages ---

func runMessageTransfer(args []string) error {
	f, err := parseMessageFlags("transfer", args, false)
	if err != nil {
		return err
	}
	var body tap.TransferBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewTransferMessage(f.from, f.to, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessagePayment(args []string) error {
	f, err := parseMessageFlags("payment", args, false)
	if err != nil {
		return err
	}
	var body tap.PaymentBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewPaymentMessage(f.from, f.to, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageRFQ(args []string) error {
	f, err := parseMessageFlags("rfq", args, false)
	if err != nil {
		return err
	}
	var body tap.RFQBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewRFQMessage(f.from, f.to, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageLock(args []string) error {
	f, err := parseMessageFlags("lock", args, false)
	if err != nil {
		return err
	}
	var body tap.LockBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewLockMessage(f.from, f.to, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageConnect(args []string) error {
	f, err := parseMessageFlags("connect", args, false)
	if err != nil {
		return err
	}
	var body tap.ConnectBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewConnectMessage(f.from, f.to, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

// --- Reply messages ---

func runMessageAuthorize(args []string) error {
	f, err := parseMessageFlags("authorize", args, true)
	if err != nil {
		return err
	}
	var body tap.AuthorizeBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewAuthorizeMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageAuthorizationRequired(args []string) error {
	f, err := parseMessageFlags("authorization-required", args, true)
	if err != nil {
		return err
	}
	var body tap.AuthorizationRequiredBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewAuthorizationRequiredMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageSettle(args []string) error {
	f, err := parseMessageFlags("settle", args, true)
	if err != nil {
		return err
	}
	var body tap.SettleBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewSettleMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageReject(args []string) error {
	f, err := parseMessageFlags("reject", args, true)
	if err != nil {
		return err
	}
	var body tap.RejectBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewRejectMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageCancel(args []string) error {
	f, err := parseMessageFlags("cancel", args, true)
	if err != nil {
		return err
	}
	var body tap.CancelBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewCancelMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageRevert(args []string) error {
	f, err := parseMessageFlags("revert", args, true)
	if err != nil {
		return err
	}
	var body tap.RevertBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewRevertMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageCapture(args []string) error {
	f, err := parseMessageFlags("capture", args, true)
	if err != nil {
		return err
	}
	var body tap.CaptureBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewCaptureMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageQuote(args []string) error {
	f, err := parseMessageFlags("quote", args, true)
	if err != nil {
		return err
	}
	var body tap.QuoteBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewQuoteMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageAddAgents(args []string) error {
	f, err := parseMessageFlags("add-agents", args, true)
	if err != nil {
		return err
	}
	var body tap.AddAgentsBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewAddAgentsMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageRemoveAgent(args []string) error {
	f, err := parseMessageFlags("remove-agent", args, true)
	if err != nil {
		return err
	}
	var body tap.RemoveAgentBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewRemoveAgentMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageReplaceAgent(args []string) error {
	f, err := parseMessageFlags("replace-agent", args, true)
	if err != nil {
		return err
	}
	var body tap.ReplaceAgentBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewReplaceAgentMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageUpdateAgent(args []string) error {
	f, err := parseMessageFlags("update-agent", args, true)
	if err != nil {
		return err
	}
	var body tap.UpdateAgentBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewUpdateAgentMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageUpdateParty(args []string) error {
	f, err := parseMessageFlags("update-party", args, true)
	if err != nil {
		return err
	}
	var body tap.UpdatePartyBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewUpdatePartyMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageUpdatePolicies(args []string) error {
	f, err := parseMessageFlags("update-policies", args, true)
	if err != nil {
		return err
	}
	var body tap.UpdatePoliciesBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewUpdatePoliciesMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}

func runMessageConfirmRelationship(args []string) error {
	f, err := parseMessageFlags("confirm-relationship", args, true)
	if err != nil {
		return err
	}
	var body tap.ConfirmRelationshipBody
	if err := json.Unmarshal(f.body, &body); err != nil {
		return fmt.Errorf("parse body JSON: %w", err)
	}
	msg, err := tap.NewConfirmRelationshipMessage(f.from, f.to, f.thid, &body)
	if err != nil {
		return err
	}
	return writeMessage(msg)
}
