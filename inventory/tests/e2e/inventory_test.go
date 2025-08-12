package e2e

import (
	"context"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient invV1.InventoryServiceClient
	)
	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		Expect(err).ToNot(HaveOccurred(), "expected to connect to gRPC app")

		inventoryClient = invV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		err := env.ClearInventoryCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "expected to clear collection inventory")

		cancel()
	})

	Describe("GetPart", func() {
		var partUuid string

		BeforeEach(func() {
			var err error

			partUuid, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "expected successful insertion in MongoDB")
		})

		It("need to return part via Uuid", func() {
			resp, err := inventoryClient.GetPart(ctx, &invV1.GetPartRequest{
				Uuid: partUuid,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().GetUuid()).To(Equal(partUuid))
			Expect(resp.GetPart().GetManufacturer()).ToNot(BeNil())
			Expect(resp.GetPart().GetTags()).ToNot(BeNil())
			Expect(resp.GetPart().GetDimensions()).ToNot(BeNil())
		})

		It("need to return not found error", func() {
			resp, err := inventoryClient.GetPart(ctx, &invV1.GetPartRequest{
				Uuid: "what is that",
			})
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})
	Describe("ListParts", func() {
		var partUuids []string
		BeforeEach(func() {
			var err error

			partUuids, err = env.InsertTestParts(ctx)

			Expect(err).ToNot(HaveOccurred(), "expected successful multiple insertion in MongoDB")
			Expect(err).To(BeNil())

		})

		It("need to filter parts via uuids", func() {
			resp, err := inventoryClient.ListParts(ctx, &invV1.ListPartsRequest{
				Filter: &invV1.PartsFilter{
					Uuids: partUuids,
				},
			})

			// logs, logErr := env.App.Logs(ctx)
			// if logErr == nil {
			// 	defer logs.Close()
			// 	logBytes, _ := io.ReadAll(logs)
			// 	fmt.Printf("=== CONTAINER LOGS ===\n%s\n=== END LOGS ===\n", string(logBytes))
			// }

			Expect(err).ToNot(HaveOccurred())
			Expect(resp).ToNot(BeNil())
			Expect(resp.GetParts()).To(HaveLen(2))
		})
		It("needs to return empty res", func() {
			resp, err := inventoryClient.ListParts(ctx, &invV1.ListPartsRequest{
				Filter: &invV1.PartsFilter{
					Uuids: []string{"govnooo"},
				},
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).To(BeNil())
		})

	})
})
