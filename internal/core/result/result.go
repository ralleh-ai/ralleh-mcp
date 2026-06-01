package result

import "time"

// Capabilities is returned by services so clients can see hard safety limits.
type Capabilities struct {
	CanSearch             bool `json:"canSearch"`
	CanCompare            bool `json:"canCompare"`
	CanVerify             bool `json:"canVerify"`
	CanPurchase           bool `json:"canPurchase"`
	CanBook               bool `json:"canBook"`
	CanUseCreditCard      bool `json:"canUseCreditCard"`
	CanEnterPassengerInfo bool `json:"canEnterPassengerInfo"`
}

// SourceDiagnostic records what happened for one source during a bounded run.
type SourceDiagnostic struct {
	SourceID    string        `json:"sourceId"`
	Status      string        `json:"status"`
	Mode        string        `json:"mode,omitempty"`
	Duration    time.Duration `json:"duration"`
	ErrorType   string        `json:"errorType,omitempty"`
	Error       string        `json:"error,omitempty"`
	ResultCount int           `json:"resultCount"`
}

// SourcePlan is the executed plan, not merely what the LLM requested.
type SourcePlan struct {
	Collection       string             `json:"collection"`
	RequestedSources []string           `json:"requestedSources,omitempty"`
	AcceptedSources  []string           `json:"acceptedSources"`
	RejectedSources  []string           `json:"rejectedSources,omitempty"`
	AddedSources     []string           `json:"addedSources,omitempty"`
	BudgetProfile    string             `json:"budgetProfile"`
	Diagnostics      []SourceDiagnostic `json:"diagnostics,omitempty"`
}

// AffiliateDisclosure tells the client whether a one-time disclosure is needed.
type AffiliateDisclosure struct {
	Required bool   `json:"required"`
	Text     string `json:"text,omitempty"`
}
