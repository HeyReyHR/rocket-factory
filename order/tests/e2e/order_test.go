package e2e

import (
	"context"
	"fmt"

	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OrderService", func() {
	var (
		ctx         context.Context
		cancel      context.CancelFunc
		orderClient *orderV1.Client
		orderUuid   string
	)
	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		var err error
		orderClient, err = orderV1.NewClient("http://" + env.App.Address())
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		cancel()
	})
	Describe("PostOrder", func() {
		It("need to create order", func() {
			partUuid, _, err := env.InsertTestParts(ctx) // idk maybe ts is inconsistent (f.e. PostOrder would run faster than insertion) and overall looks like ass yk
			if err != nil {
				Fail(fmt.Sprintf("Cannot insert test parts: %v", err))
			}

			resp, err := orderClient.PostOrder(ctx, &orderV1.CreateOrderRequest{
				UserUUID:  "1",
				PartUuids: []string{partUuid},
			})
			Expect(err).ToNot(HaveOccurred())

			if createResp, ok := resp.(*orderV1.CreateOrderResponse); ok {
				orderUuid = createResp.UUID
				Expect(orderUuid).ToNot(BeEmpty())
			} else {
				Fail(fmt.Sprintf("Response body is not proper and contains error: %v", createResp))
			}
		})

		It("need return out of stock error", func() {
			_, partUuidOutOfStock, err := env.InsertTestParts(ctx) // same as line 34
			if err != nil {
				Fail(fmt.Sprintf("Cannot insert test parts: %v", err))
			}
			resp, err := orderClient.PostOrder(ctx, &orderV1.CreateOrderRequest{
				UserUUID:  "1",
				PartUuids: []string{partUuidOutOfStock},
			})
			Expect(err).ToNot(HaveOccurred())

			if createResp, ok := resp.(*orderV1.BadRequestError); ok {
				Expect(createResp.Code).To(Equal(400))
			} else {
				Fail(fmt.Sprintf("Response body is not proper and doesnt contain expected error: %v", createResp))
			}
		})
		It("need return parts not found error", func() {
			resp, err := orderClient.PostOrder(ctx, &orderV1.CreateOrderRequest{
				UserUUID:  "1",
				PartUuids: []string{"-"},
			})
			Expect(err).ToNot(HaveOccurred())

			if createResp, ok := resp.(*orderV1.BadRequestError); ok {
				Expect(createResp.Code).To(Equal(400))
			} else {
				Fail(fmt.Sprintf("Response body is not proper and doesnt contain expected error: %v", createResp))
			}
		})
	})

	Describe("GetOrder", func() {
		It("need to get order", func() {
			resp, err := orderClient.GetOrder(ctx, orderV1.GetOrderParams{
				OrderUUID: orderUuid,
			})
			Expect(err).ToNot(HaveOccurred())

			if order, ok := resp.(*orderV1.OrderDto); ok {
				Expect(order.UUID).To(Equal(orderUuid))
				Expect(order.UserUUID).To(Equal("1"))
			} else {
				Fail(fmt.Sprintf("Response body is not proper and contains error: %v", resp))
			}
		})

		It("needs to return 404 error", func() {
			resp, err := orderClient.GetOrder(ctx, orderV1.GetOrderParams{
				OrderUUID: "-",
			})
			Expect(err).ToNot(HaveOccurred())

			if order, ok := resp.(*orderV1.NotFoundError); ok {
				Expect(order.Code).To(Equal(404))
			} else {
				Fail(fmt.Sprintf("Response body is not proper and doesnt contain expected error: %v", resp))
			}
		})
	})

	Describe("Pay", func() {
		It("needs to pay order", func() {
			payResp, err := orderClient.PayOrder(ctx, &orderV1.OrderPayRequest{
				PaymentMethod: orderV1.PaymentMethodCARD,
			}, orderV1.PayOrderParams{
				OrderUUID: orderUuid,
			})
			Expect(err).ToNot(HaveOccurred())

			if resp, ok := payResp.(*orderV1.OrderPayResponse); ok {
				Expect(resp.TransactionUUID).ToNot(BeNil())
				Expect(resp.TransactionUUID).ToNot(Equal(""))
				order, err := orderClient.GetOrder(ctx, orderV1.GetOrderParams{
					OrderUUID: orderUuid,
				})
				Expect(err).ToNot(HaveOccurred())
				if order, ok := order.(*orderV1.OrderDto); ok {
					Expect(order.PaymentMethod).To(Equal(orderV1.OptPaymentMethod{
						Value: orderV1.PaymentMethodCARD,
						Set:   true,
					}))
				}
			} else {
				Fail(fmt.Sprintf("Response body is not proper and contains error: %v", resp))
			}
		})
		It("needs to return already paid error", func() {
			resp, err := orderClient.PayOrder(ctx, &orderV1.OrderPayRequest{
				PaymentMethod: orderV1.PaymentMethodCARD,
			}, orderV1.PayOrderParams{
				OrderUUID: orderUuid,
			})
			Expect(err).ToNot(HaveOccurred())
			if badReqResp, ok := resp.(*orderV1.BadRequestError); ok {
				Expect(badReqResp.Code).To(Equal(400)) // TODO MESSAGE INSTEAD OF CODE
			} else {
				Fail(fmt.Sprintf("Response body is not proper and doesnt contain expected error: %v", resp))
			}
		})
		It("needs to return already cancelled", func() {
			orderUuidCancelled, err := env.InsertCancelledOrder(ctx) // same as line 34
			if err != nil {
				Fail(fmt.Sprintf("Cannot insert order: %v", err))
			}

			resp, err := orderClient.PayOrder(ctx, &orderV1.OrderPayRequest{
				PaymentMethod: orderV1.PaymentMethodCARD,
			}, orderV1.PayOrderParams{
				OrderUUID: orderUuidCancelled,
			})
			Expect(err).ToNot(HaveOccurred())

			if badReqResp, ok := resp.(*orderV1.BadRequestError); ok {
				Expect(badReqResp.Code).To(Equal(400))
			} else {
				Fail(fmt.Sprintf("Response body is not proper and doesnt contain expected error: %v", resp))
			}
		})
		It("needs to return not found error", func() {
			resp, err := orderClient.PayOrder(ctx, &orderV1.OrderPayRequest{
				PaymentMethod: orderV1.PaymentMethodCARD,
			}, orderV1.PayOrderParams{
				OrderUUID: "-",
			})
			Expect(err).ToNot(HaveOccurred())

			if resp, ok := resp.(*orderV1.NotFoundError); !ok {
				Fail(fmt.Sprintf("Response body is not proper and doesnt contain expected error: %v", resp))
			}
		})
	})
	Describe("Cancel", func() {
		It("needs to cancel order", func() {
			orderUuidPendingPayment, err := env.InsertOrder(ctx) // same as line 34
			if err != nil {
				Fail(fmt.Sprintf("Cannot insert order: %v", err))
			}

			resp, err := orderClient.CancelOrder(ctx, orderV1.CancelOrderParams{
				OrderUUID: orderUuidPendingPayment,
			})
			Expect(err).ToNot(HaveOccurred())
			if resp, ok := resp.(*orderV1.CancelOrderNoContent); !ok {
				Fail(fmt.Sprintf("Response body is not proper and contains error: %v", resp))
			}
		})
		It("needs to return already paid error", func() {
			resp, err := orderClient.CancelOrder(ctx, orderV1.CancelOrderParams{
				OrderUUID: orderUuid,
			})
			Expect(err).ToNot(HaveOccurred())
			if resp, ok := resp.(*orderV1.ConfilctError); !ok {
				Fail(fmt.Sprintf("Response body is not proper and doesnt contain expected error: %v", resp))
			}
		})
		It("needs to return not found error", func() {
			resp, err := orderClient.CancelOrder(ctx, orderV1.CancelOrderParams{
				OrderUUID: "-",
			})
			Expect(err).ToNot(HaveOccurred())

			if resp, ok := resp.(*orderV1.NotFoundError); !ok {
				Fail(fmt.Sprintf("Response body is not proper and doesnt contain expected error: %v", resp))
			}
		})
	})
})
