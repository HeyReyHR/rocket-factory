package telegram

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"
	"time"

	"github.com/HeyReyHR/rocket-factory/notification/internal/client/http"
	"github.com/HeyReyHR/rocket-factory/notification/internal/model"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

const chatID = 1905584478

//go:embed templates/*.tmpl
var templateFS embed.FS

type orderPaidTemplateData struct {
	OrderUuid       string
	UserUuid        string
	PaymentMethod   string
	TransactionUuid string
}

type orderAssembledTemplateData struct {
	OrderUuid    string
	UserUuid     string
	BuildTimeSec int64
}

var orderAssembledTemplate = template.Must(template.ParseFS(templateFS, "templates/order_assembled_notification.tmpl"))

var orderPaidTemplate = template.Must(template.ParseFS(templateFS, "templates/order_paid_notification.tmpl"))

type service struct {
	telegramClient http.TelegramClient
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, order model.OrderPaidEvent) error {
	message, err := s.buildOrderPaidMessage(order)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) SendOrderAssembledNotification(ctx context.Context, order model.OrderAssembledEvent) error {
	message, err := s.buildOrderAssembledMessage(order)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) buildOrderPaidMessage(order model.OrderPaidEvent) (string, error) {
	fmt.Println(order.PaymentMethod)
	data := orderPaidTemplateData{
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PaymentMethod:   order.PaymentMethod,
		TransactionUuid: order.TransactionUuid,
	}

	var buf bytes.Buffer
	err := orderPaidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) buildOrderAssembledMessage(order model.OrderAssembledEvent) (string, error) {
	buildTime := time.Duration(order.BuildTimeSec) / time.Second
	data := orderAssembledTemplateData{
		OrderUuid:    order.OrderUuid,
		UserUuid:     order.UserUuid,
		BuildTimeSec: int64(buildTime),
	}

	var buf bytes.Buffer
	err := orderAssembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
