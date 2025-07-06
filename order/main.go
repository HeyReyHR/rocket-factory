package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	paymentServiceAddress   = "localhost:50052"
	inventoryServiceAddress = "localhost:50051"
	httpPort                = "8080"
	readHeaderTimeout       = 5 * time.Second
	shutdownTimeout         = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}
type OrderHandler struct {
	storage   *OrderStorage
	inventory invV1.InventoryServiceClient
	payment   payV1.PaymentServiceClient
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func NewOrderHandler(storage *OrderStorage, inventory invV1.InventoryServiceClient, payment payV1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:   storage,
		inventory: inventory,
		payment:   payment,
	}
}

func convertPayment(method orderV1.PaymentMethod) payV1.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCARD:
		return payV1.PaymentMethod_CARD
	case orderV1.PaymentMethodCREDITCARD:
		return payV1.PaymentMethod_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return payV1.PaymentMethod_INVESTOR_MONEY
	case orderV1.PaymentMethodSBP:
		return payV1.PaymentMethod_SBP
	default:
		return payV1.PaymentMethod_UNKNOWN
	}
}

func (s *OrderStorage) UpdateOrder(orderUUID string, order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[orderUUID] = order
}

func (s *OrderStorage) GetOrder(orderUUID string) *orderV1.OrderDto {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderStorage) CreateOrder(orderUUID, userUUID string, partsUUIDs []string, totalPrice float64) *orderV1.OrderDto {
	s.mu.Lock()
	defer s.mu.Unlock()

	order := &orderV1.OrderDto{
		UserUUID:   userUUID,
		PartUuids:  partsUUIDs,
		UUID:       orderUUID,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
		TotalPrice: totalPrice,
	}

	s.orders[orderUUID] = order

	return order
}

func (h *OrderHandler) GetOrder(_ context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order By UUID " + params.OrderUUID + " not found",
		}, nil
	}
	return order, nil
}

func (h *OrderHandler) PostOrder(ctx context.Context, r *orderV1.CreateOrderRequest) (orderV1.PostOrderRes, error) {
	resp, err := h.inventory.ListParts(ctx, &invV1.ListPartsRequest{Filter: &invV1.PartsFilter{
		Uuids: r.PartUuids,
	}})
	if err != nil {
		return nil, err
	}

	if len(resp.Parts) != len(r.PartUuids) {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Some parts haven't been found",
		}, nil
	}

	newUUID := uuid.NewString()

	var totalPrice float64
	for _, part := range resp.Parts {
		totalPrice += part.Price
	}

	h.storage.CreateOrder(newUUID, r.UserUUID, r.PartUuids, totalPrice)
	return &orderV1.CreateOrderResponse{
		UUID:       newUUID,
		TotalPrice: totalPrice,
	}, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, r *orderV1.OrderPayRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order By UUID " + params.OrderUUID + " not found",
		}, nil
	}
	payment, err := h.payment.PayOrder(ctx, &payV1.PayOrderRequest{
		PaymentMethod: convertPayment(r.PaymentMethod),
		OrderUuid:     order.UUID,
		UserUuid:      order.UserUUID,
	})
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Could not proceed payment",
		}, nil
	}

	order.TransactionUUID = orderV1.OptString{
		Value: payment.TransactionUuid,
		Set:   true,
	}
	order.PaymentMethod = orderV1.OptPaymentMethod{
		Value: r.PaymentMethod,
		Set:   true,
	}
	order.Status = orderV1.OrderStatusPAID

	h.storage.UpdateOrder(params.OrderUUID, order)

	return &orderV1.OrderPayResponse{
		TransactionUUID: payment.TransactionUuid,
	}, nil
}

func (h *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order By UUID " + params.OrderUUID + " not found",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConfilctError{
			Code:    409,
			Message: "Order By UUID " + params.OrderUUID + " already paid, cannot be cancelled",
		}, nil
	}

	order.Status = orderV1.OrderStatusCANCELLED

	h.storage.UpdateOrder(params.OrderUUID, order)

	return &orderV1.CancelOrderNoContent{}, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.InternalServerErrorStatusCode {
	return &orderV1.InternalServerErrorStatusCode{
		StatusCode: 500,
		Response: orderV1.InternalServerError{
			Code:    500,
			Message: "Internal Server Error: " + err.Error(),
		},
	}
}

func main() {
	connPay, err := grpc.NewClient(
		paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to payment service: %s", err)
	}

	defer func() {
		if cerr := connPay.Close(); cerr != nil {
			log.Printf("failed to close connection with payment: %s", cerr)
		}
	}()

	payment := payV1.NewPaymentServiceClient(connPay)

	connInv, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to inventory service: %s", err)
	}

	defer func() {
		if cerr := connInv.Close(); cerr != nil {
			log.Printf("failed to close connection with inventory: %s", cerr)
		}
	}()

	inventory := invV1.NewInventoryServiceClient(connInv)

	storage := NewOrderStorage()

	orderHandler := NewOrderHandler(storage, inventory, payment)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("Error occured when creating OpenAPI server: %s", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 Starting server on port %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Error occurred when starting server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Error occurred when shutting down server: %s\n", err)
	}

	log.Println("✅ Server stopped")
}
