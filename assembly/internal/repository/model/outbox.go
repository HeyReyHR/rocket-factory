package model

type OrderAssembled struct {
	EventUuid string
	EventType EventType
	Payload   OrderAssembledPayload
	Status    Status
}

type OrderAssembledPayload struct {
	OrderUuid    string `json:"order_uuid"`
	UserUuid     string `json:"user_uuid"`
	BuildTimeSec int64  `json:"build_time_sec"`
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
