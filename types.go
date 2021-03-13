package types

// CustomerInitiation s
type CustomerInitiation struct {
	Initiation Initiation
	Customer   CustomerReference
}

// Initiation s
type Initiation struct {
	GroupHeader         GroupHeader
	PaymentInstructions PaymentInstructions
}

// GroupHeader s
type GroupHeader struct {
	MessageID            string
	CreationDateTime     string
	NumberOfTransactions string
	ControlSum           string
	InitiatingParty      InitiatingPartyReference
}

// PaymentInstruction s
type PaymentInstruction struct {
	PaymentID              string
	NumberOfTransactions   string
	ControlSum             string
	RequestedExecutionDate string
	ExpectedExecutionDate  string
	Debtor                 AccountReference
}

// PaymentInstructions a
type PaymentInstructions []PaymentInstruction

// AccountReference s
type AccountReference struct {
	Name string
	IBAN string
	BIC  string
}

// InitiatingPartyReference s
type InitiatingPartyReference struct {
	InitiatingPartyID string
}

// UserReference s
type UserReference struct {
	Email string
}

// CustomerReference s
type CustomerReference struct {
	CustomerID      string
	Name            string
	InitiatingParty InitiatingPartyReference
}
