package app

import (
	"context"
	"fmt"

	"github.com/HeyReyHR/rocket-factory/assembly/internal/config"
	"github.com/HeyReyHR/rocket-factory/assembly/internal/converter/kafka/decoder"
	"github.com/HeyReyHR/rocket-factory/assembly/internal/service"
	orderConsumer "github.com/HeyReyHR/rocket-factory/assembly/internal/service/consumer/assembly_consumer"
	shipProducer "github.com/HeyReyHR/rocket-factory/assembly/internal/service/producer/assembly_producer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/closer"
	wrappedKafkaConsumer "github.com/HeyReyHR/rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/HeyReyHR/rocket-factory/platform/pkg/kafka/producer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/kafka"
	"github.com/IBM/sarama"

	kafkaConverter "github.com/HeyReyHR/rocket-factory/assembly/internal/converter/kafka"
	wrappedKafka "github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
)

type diContainer struct {
	shipProducerService  service.ShipProducerService
	orderConsumerService service.ConsumerService

	consumerGroup     sarama.ConsumerGroup
	orderPaidConsumer wrappedKafka.Consumer

	orderPaidDecoder      kafkaConverter.OrderPaidDecoder
	syncProducer          sarama.SyncProducer
	shipAssembledProducer wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) ShipAssembledProducerService() service.ShipProducerService {
	if d.shipProducerService == nil {
		d.shipProducerService = shipProducer.NewService(d.shipAssembledProducer)
	}

	return d.shipProducerService
}

func (d *diContainer) OrderConsumerService() service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(d.OrderPaidConsumer(), d.OrderPaidDecoder())
	}

	return d.orderConsumerService
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().ShipAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) ShipAssembledProducer() wrappedKafka.Producer {
	if d.shipAssembledProducer == nil {
		d.shipAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().ShipAssembledProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.shipAssembledProducer
}
