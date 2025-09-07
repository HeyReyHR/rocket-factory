package decoder

import (
	"fmt"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	eventsV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type decoder struct{}

func NewShipAssembledDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUuid:    pb.EventUuid,
		OrderUuid:    pb.OrderUuid,
		UserUuid:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
