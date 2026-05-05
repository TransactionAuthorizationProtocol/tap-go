package tap

import (
	"strings"
	"testing"
)

func TestTypeConstants(t *testing.T) {
	types := map[string]string{
		"Transfer":              TypeTransfer,
		"Payment":               TypePayment,
		"RFQ":                   TypeRFQ,
		"Quote":                 TypeQuote,
		"Lock":                  TypeLock,
		"Authorize":             TypeAuthorize,
		"AuthorizationRequired": TypeAuthorizationRequired,
		"Settle":                TypeSettle,
		"Reject":                TypeReject,
		"Cancel":                TypeCancel,
		"Revert":                TypeRevert,
		"Capture":               TypeCapture,
		"UpdateAgent":           TypeUpdateAgent,
		"UpdateParty":           TypeUpdateParty,
		"AddAgents":             TypeAddAgents,
		"ReplaceAgent":          TypeReplaceAgent,
		"RemoveAgent":           TypeRemoveAgent,
		"ConfirmRelationship":   TypeConfirmRelationship,
		"UpdatePolicies":        TypeUpdatePolicies,
		"Connect":               TypeConnect,
	}

	for name, typ := range types {
		if !strings.HasPrefix(typ, "https://tap.rsvp/schema/1.0#") {
			t.Errorf("%s: type %q does not have expected prefix", name, typ)
		}
		if !strings.HasSuffix(typ, name) {
			t.Errorf("%s: type %q does not end with %q", name, typ, name)
		}
	}
}

func TestTAPContext(t *testing.T) {
	if TAPContext != "https://tap.rsvp/schema/1.0" {
		t.Errorf("TAPContext: got %q", TAPContext)
	}
}
