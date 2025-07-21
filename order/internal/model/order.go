package model

type Order struct {
	Uuid            string
	UserUuid        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUuid *string
	PaymentMethod   *PaymentMethod
	Status          Status
}

type PaymentMethod int

const (
	UNKNOWN_METHOD PaymentMethod = iota
	CARD
	SBP
	CREDIT_CARD
	INVESTOR_MONEY
)

type Status int

const (
	UNKNOWN_STATUS Status = iota
	PENDING_PAYMENT
	PAID
	CANCELLED
)
