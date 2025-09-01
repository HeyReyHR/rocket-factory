package model

type OrderPaidEvent struct {
	EventUuid       string
	OrderUuid       string
	UserUuid        string
	PaymentMethod   string
	TransactionUuid string
}

type OrderAssembledEvent struct {
	EventUuid    string
	OrderUuid    string
	UserUuid     string
	BuildTimeSec int64
}

type Status string

const (
	UnknownStatus Status = "unknown_status"
	Done          Status = "done"
	PendingStatus Status = "pending"
)

type EventType string

const (
	UnknownType             EventType = "unknown_type"
	OrderAssembledEventType EventType = "order_assembled"
)
