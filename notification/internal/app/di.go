package app

import (
	"context"
	"fmt"

	httpClient "github.com/HeyReyHR/rocket-factory/notification/internal/client/http"
	telegramClient "github.com/HeyReyHR/rocket-factory/notification/internal/client/http/telegram"
	"github.com/HeyReyHR/rocket-factory/notification/internal/config"
	kafkaConverter "github.com/HeyReyHR/rocket-factory/notification/internal/converter/kafka"
	"github.com/HeyReyHR/rocket-factory/notification/internal/converter/kafka/decoder"
	"github.com/HeyReyHR/rocket-factory/notification/internal/service"
	orderAssembledConsumer "github.com/HeyReyHR/rocket-factory/notification/internal/service/consumer/order_assembled_consumer"
	orderPaidConsumer "github.com/HeyReyHR/rocket-factory/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/HeyReyHR/rocket-factory/notification/internal/service/telegram"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/HeyReyHR/rocket-factory/platform/pkg/kafka/consumer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/kafka"
	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"
)

type diContainer struct {
	orderPaidConsumerService      service.OrderPaidConsumerService
	orderAssembledConsumerService service.OrderAssembledConsumerService

	orderPaidConsumerGroup      sarama.ConsumerGroup
	orderAssembledConsumerGroup sarama.ConsumerGroup

	orderPaidConsumer wrappedKafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder

	orderAssembledConsumer wrappedKafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder

	telegramService service.TelegramService
	telegramClient  httpClient.TelegramClient
	telegramBot     *bot.Bot
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.OrderAssembledConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = orderAssembledConsumer.NewService(d.OrderAssembledConsumer(), d.OrderAssembledDecoder(), d.TelegramService(ctx))
	}

	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumer.NewService(d.OrderPaidConsumer(), d.OrderPaidDecoder(), d.TelegramService(ctx))
	}

	return d.orderPaidConsumerService
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderPaidConsumerGroup.Close()
		})

		d.orderPaidConsumerGroup = consumerGroup
	}

	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderAssembledConsumerGroup.Close()
		})

		d.orderAssembledConsumerGroup = consumerGroup
	}

	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(
			d.TelegramClient(ctx),
		)
	}
	return d.telegramService
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().Telegram.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}
