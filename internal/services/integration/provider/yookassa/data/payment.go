package data

const (
	PaymentStatusPending   = "pending"
	PaymentStatusSucceeded = "succeeded"
	PaymentStatusCanceled  = "canceled"
)

type PaymentAmount struct {
	Currency string
	Amount   float64
}

type PaymentMethod struct {
	Type string
}

const (
	PaymentConfirmationTypeRedirect = "redirect"
)

type PaymentConfirmation struct {
	Type      string
	ReturnURL string
	Locale    string
}

type CreatePaymentRequest struct {
	Confirmation  PaymentConfirmation
	PaymentMethod PaymentMethod
	Amount        PaymentAmount
	Description   string
	Capture       bool
}

type CreatePaymentResponse struct {
	ID              string
	ConfirmationURL string
	Status          string
}
