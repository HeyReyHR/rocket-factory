package decoder

import (
	"fmt"

	"github.com/HeyReyHR/rocket-factory/notification/internal/model"
	eventsV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

func NewOrderAssembledDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) DecodeOrderAssembled(data []byte) (model.OrderAssembledEvent, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderAssembledEvent{
		EventUuid:    pb.EventUuid,
		OrderUuid:    pb.OrderUuid,
		UserUuid:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
