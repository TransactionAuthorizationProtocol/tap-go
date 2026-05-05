package tap

const (
	// TAPContext is the JSON-LD context for all TAP messages.
	TAPContext = "https://tap.rsvp/schema/1.0"

	// Transaction message types
	TypeTransfer = "https://tap.rsvp/schema/1.0#Transfer"
	TypePayment  = "https://tap.rsvp/schema/1.0#Payment"
	TypeRFQ      = "https://tap.rsvp/schema/1.0#RFQ"
	TypeQuote    = "https://tap.rsvp/schema/1.0#Quote"
	TypeLock     = "https://tap.rsvp/schema/1.0#Lock"

	// Authorization flow message types
	TypeAuthorize             = "https://tap.rsvp/schema/1.0#Authorize"
	TypeAuthorizationRequired = "https://tap.rsvp/schema/1.0#AuthorizationRequired"
	TypeSettle                = "https://tap.rsvp/schema/1.0#Settle"
	TypeReject                = "https://tap.rsvp/schema/1.0#Reject"
	TypeCancel                = "https://tap.rsvp/schema/1.0#Cancel"
	TypeRevert                = "https://tap.rsvp/schema/1.0#Revert"
	TypeCapture               = "https://tap.rsvp/schema/1.0#Capture"

	// Agent management message types
	TypeUpdateAgent  = "https://tap.rsvp/schema/1.0#UpdateAgent"
	TypeUpdateParty  = "https://tap.rsvp/schema/1.0#UpdateParty"
	TypeAddAgents    = "https://tap.rsvp/schema/1.0#AddAgents"
	TypeReplaceAgent = "https://tap.rsvp/schema/1.0#ReplaceAgent"
	TypeRemoveAgent  = "https://tap.rsvp/schema/1.0#RemoveAgent"

	// Relationship message types
	TypeConfirmRelationship = "https://tap.rsvp/schema/1.0#ConfirmRelationship"
	TypeUpdatePolicies      = "https://tap.rsvp/schema/1.0#UpdatePolicies"
	TypeConnect             = "https://tap.rsvp/schema/1.0#Connect"
)

// AllTypes returns all known TAP message types.
func AllTypes() []string {
	return []string{
		TypeTransfer,
		TypePayment,
		TypeRFQ,
		TypeQuote,
		TypeLock,
		TypeAuthorize,
		TypeAuthorizationRequired,
		TypeSettle,
		TypeReject,
		TypeCancel,
		TypeRevert,
		TypeCapture,
		TypeUpdateAgent,
		TypeUpdateParty,
		TypeAddAgents,
		TypeReplaceAgent,
		TypeRemoveAgent,
		TypeConfirmRelationship,
		TypeUpdatePolicies,
		TypeConnect,
	}
}
