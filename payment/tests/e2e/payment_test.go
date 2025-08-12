package e2e

import (
	"context"
	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ = Describe("PaymentService", func() {
	var (
		ctx           context.Context
		cancel        context.CancelFunc
		paymentClient payV1.PaymentServiceClient
	)
	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		Expect(err).ToNot(HaveOccurred(), "expected to connect to gRPC app")

		paymentClient = payV1.NewPaymentServiceClient(conn)
	})
	AfterEach(func() {
		cancel()
	})

	Describe("PayOrder", func() {
		orderUuid := uuid.NewString()
		It("need to return transaction uuid", func() {
			resp, err := paymentClient.PayOrder(ctx, &payV1.PayOrderRequest{
				OrderUuid:     orderUuid,
				PaymentMethod: payV1.PaymentMethod_SBP,
			})
			//logs, logErr := env.App.Logs(ctx)
			//if logErr == nil {
			//	defer logs.Close()
			//	logBytes, _ := io.ReadAll(logs)
			//	fmt.Printf("=== CONTAINER LOGS ===\n%s\n=== END LOGS ===\n", string(logBytes))
			//}
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetTransactionUuid()).ToNot(BeNil())
		})

		It("need to return error when unknown payment method", func() {
			resp, err := paymentClient.PayOrder(ctx, &payV1.PayOrderRequest{
				OrderUuid:     orderUuid,
				PaymentMethod: payV1.PaymentMethod_UNKNOWN,
			})
			Expect(err).To(HaveOccurred())
			Expect(resp.GetTransactionUuid()).To(Equal(""))
		})

		It("need to return error when nil order uuid", func() {
			resp, err := paymentClient.PayOrder(ctx, &payV1.PayOrderRequest{
				OrderUuid:     "",
				PaymentMethod: payV1.PaymentMethod_SBP,
			})
			Expect(err).To(HaveOccurred())
			Expect(resp.GetTransactionUuid()).To(Equal(""))
		})
	})
})
