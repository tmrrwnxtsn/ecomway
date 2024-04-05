package model

type TransactionType string

const (
	TransactionTypePayment TransactionType = "payment"
	TransactionTypePayout  TransactionType = "payout"
)
