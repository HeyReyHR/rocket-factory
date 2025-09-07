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

type PaymentMethod string

const (
	UNKNOWN_METHOD PaymentMethod = "unknown_method"
	CARD           PaymentMethod = "card"
	SBP            PaymentMethod = "sbp"
	CREDIT_CARD    PaymentMethod = "credit_card"
	INVESTOR_MONEY PaymentMethod = "investor_money"
)

type Status string

const (
	UNKNOWN_STATUS  Status = "unknown_status"
	PENDING_PAYMENT Status = "pending_payment"
	PAID            Status = "paid"
	CANCELLED       Status = "cancelled"
	ASSEMBLED       Status = "assembled"
)
