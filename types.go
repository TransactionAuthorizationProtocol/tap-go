package tap

import "encoding/json"

// Party represents a real-world entity (legal or natural person) involved in a transaction.
type Party struct {
	ID          string           `json:"@id"`
	Type        string           `json:"@type,omitempty"`
	Name        string           `json:"name,omitempty"`
	NameHash    string           `json:"nameHash,omitempty"`
	LEICode     string           `json:"lei:leiCode,omitempty"`
	IVMS101     *json.RawMessage `json:"ivms101,omitempty"`
	MCC         string           `json:"mcc,omitempty"`
	Account     string           `json:"account,omitempty"`
	URL         string           `json:"url,omitempty"`
	Logo        string           `json:"logo,omitempty"`
	Description string           `json:"description,omitempty"`
	Email       string           `json:"email,omitempty"`
	Telephone   string           `json:"telephone,omitempty"`
}

// Agent represents software acting on behalf of a participant.
type Agent struct {
	ID          string   `json:"@id"`
	Type        string   `json:"@type,omitempty"`
	Role        string   `json:"role,omitempty"`
	For         ForField `json:"for,omitempty"`
	Name        string   `json:"name,omitempty"`
	NameHash    string   `json:"nameHash,omitempty"`
	LEICode     string   `json:"lei:leiCode,omitempty"`
	Policies    []Policy `json:"policies,omitempty"`
	URL         string   `json:"url,omitempty"`
	Logo        string   `json:"logo,omitempty"`
	Description string   `json:"description,omitempty"`
	Email       string   `json:"email,omitempty"`
	Telephone   string   `json:"telephone,omitempty"`
	ServiceURL  string   `json:"serviceUrl,omitempty"`
}

// ForField represents the "for" field on an Agent, which can be a single DID string
// or an array of DID strings.
type ForField struct {
	values []string
}

// NewForField creates a ForField from one or more DID strings.
func NewForField(dids ...string) ForField {
	return ForField{values: dids}
}

// Values returns the DID strings in the ForField.
func (f ForField) Values() []string {
	return f.values
}

// String returns the first DID or empty string.
func (f ForField) String() string {
	if len(f.values) > 0 {
		return f.values[0]
	}
	return ""
}

// IsEmpty returns true if the ForField has no values.
func (f ForField) IsEmpty() bool {
	return len(f.values) == 0
}

// MarshalJSON marshals the ForField as a string (single value) or array (multiple values).
func (f ForField) MarshalJSON() ([]byte, error) {
	if len(f.values) == 0 {
		return json.Marshal(nil)
	}
	if len(f.values) == 1 {
		return json.Marshal(f.values[0])
	}
	return json.Marshal(f.values)
}

// UnmarshalJSON unmarshals the ForField from a string or array of strings.
func (f *ForField) UnmarshalJSON(data []byte) error {
	// Try single string first
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		f.values = []string{single}
		return nil
	}

	// Try array of strings
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	f.values = arr
	return nil
}

// Policy represents a policy enforced by an agent.
type Policy struct {
	Type                   string   `json:"@type"`
	Context                any      `json:"@context,omitempty"`
	From                   string   `json:"from,omitempty"`
	FromAgent              string   `json:"fromAgent,omitempty"`
	FromRole               string   `json:"fromRole,omitempty"`
	AboutParty             string   `json:"aboutParty,omitempty"`
	Purpose                string   `json:"purpose,omitempty"`
	PresentationDefinition string   `json:"presentationDefinition,omitempty"`
	Nonce                  any      `json:"nonce,omitempty"`
	Codes                  []string `json:"codes,omitempty"`
}

// TransactionConstraints defines boundaries for transactions in a connection.
type TransactionConstraints struct {
	Purposes                   []string `json:"purposes,omitempty"`
	CategoryPurposes           []string `json:"categoryPurposes,omitempty"`
	Limits                     *Limits  `json:"limits,omitempty"`
	AllowedBeneficiaries       []Party  `json:"allowedBeneficiaries,omitempty"`
	AllowedSettlementAddresses []string `json:"allowedSettlementAddresses,omitempty"`
	AllowedAssets              []string `json:"allowedAssets,omitempty"`
}

// Limits defines financial limits for transactions.
type Limits struct {
	PerTransaction string `json:"per_transaction,omitempty"`
	PerDay         string `json:"per_day,omitempty"`
	PerWeek        string `json:"per_week,omitempty"`
	PerMonth       string `json:"per_month,omitempty"`
	PerYear        string `json:"per_year,omitempty"`
	Currency       string `json:"currency"`
}

// Invoice represents a structured invoice for payment information.
type Invoice struct {
	ID                          string              `json:"id"`
	IssueDate                   string              `json:"issueDate"`
	CurrencyCode                string              `json:"currencyCode"`
	LineItems                   []LineItem          `json:"lineItems"`
	Total                       float64             `json:"total"`
	SubTotal                    *float64            `json:"subTotal,omitempty"`
	TaxTotal                    *TaxTotal           `json:"taxTotal,omitempty"`
	DueDate                     string              `json:"dueDate,omitempty"`
	Note                        string              `json:"note,omitempty"`
	PaymentTerms                string              `json:"paymentTerms,omitempty"`
	AccountingCost              string              `json:"accountingCost,omitempty"`
	OrderReference              *OrderReference     `json:"orderReference,omitempty"`
	AdditionalDocumentReference []DocumentReference `json:"additionalDocumentReference,omitempty"`
}

// LineItem represents an individual item in an invoice.
type LineItem struct {
	ID          string       `json:"id"`
	Description string       `json:"description"`
	Name        string       `json:"name,omitempty"`
	Image       string       `json:"image,omitempty"`
	URL         string       `json:"url,omitempty"`
	Quantity    float64      `json:"quantity"`
	UnitCode    string       `json:"unitCode,omitempty"`
	UnitPrice   float64      `json:"unitPrice"`
	LineTotal   float64      `json:"lineTotal"`
	TaxCategory *TaxCategory `json:"taxCategory,omitempty"`
}

// TaxCategory represents tax information for a line item.
type TaxCategory struct {
	ID        string  `json:"id"`
	Percent   float64 `json:"percent"`
	TaxScheme string  `json:"taxScheme"`
}

// TaxTotal represents aggregate tax information for an invoice.
type TaxTotal struct {
	TaxAmount   float64       `json:"taxAmount"`
	TaxSubtotal []TaxSubtotal `json:"taxSubtotal,omitempty"`
}

// TaxSubtotal represents a tax breakdown by category.
type TaxSubtotal struct {
	TaxableAmount float64     `json:"taxableAmount"`
	TaxAmount     float64     `json:"taxAmount"`
	TaxCategory   TaxCategory `json:"taxCategory"`
}

// OrderReference represents information about a related order.
type OrderReference struct {
	ID        string `json:"id"`
	IssueDate string `json:"issueDate,omitempty"`
}

// DocumentReference represents a reference to an additional document.
type DocumentReference struct {
	ID           string `json:"id"`
	DocumentType string `json:"documentType,omitempty"`
	URL          string `json:"url,omitempty"`
}

// SupportedAssetPricing represents an asset with pricing information for payments.
type SupportedAssetPricing struct {
	Asset   string `json:"asset"`
	Amount  string `json:"amount"`
	Expires string `json:"expires,omitempty"`
}
