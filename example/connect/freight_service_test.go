package example

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	examplefreightv1 "github.com/einride/protoc-gen-go-aip-test/proto/gen/connect/einride/example/freight/v1"
	"github.com/einride/protoc-gen-go-aip-test/proto/gen/connect/einride/example/freight/v1/freightv1aiptest"
	"github.com/einride/protoc-gen-go-aip-test/proto/gen/connect/einride/example/freight/v1/freightv1connect"
)

// Both the Connect client and the Connect handler can be tested: the client to test a
// service over the network, the handler to test it in process.
var (
	_ freightv1aiptest.FreightService = freightv1connect.FreightServiceClient(nil)
	_ freightv1aiptest.FreightService = freightv1connect.FreightServiceHandler(nil)
)

func Test_FreightService(t *testing.T) {
	t.Parallel()
	t.Skip("this is just an example, the service is not implemented.")
	// setup server before test
	server := freightv1connect.UnimplementedFreightServiceHandler{}
	// setup test suite
	suite := freightv1aiptest.FreightServiceTestSuite{
		T:      t,
		Server: server,
	}

	// counter to keep track of unique IDs.
	var idCounter uint64

	// run tests for each resource in the service
	ctx := context.Background()
	suite.TestShipper(ctx, freightv1aiptest.FreightServiceShipperTestSuiteConfig{
		// Create should return a resource which is valid to create, i.e.
		// all required fields set.
		Create: func() *examplefreightv1.Shipper {
			return &examplefreightv1.Shipper{
				DisplayName:    "Example shipper",
				BillingAccount: "billingAccounts/12345",
			}
		},
		// IDGenerator should return a valid and unique ID to use in the Create call.
		// If non-nil, this function will be called to set the ID on all Create calls.
		// If the ID field is required, tests will fail if this is nil.
		IDGenerator: func() string {
			id := atomic.AddUint64(&idCounter, 1)
			return fmt.Sprintf("valid-id-%d", id)
		},
		// Update should return a resource which is valid to update, i.e.
		// all required fields set.
		Update: func() *examplefreightv1.Shipper {
			return &examplefreightv1.Shipper{
				DisplayName:    "Updated example shipper",
				BillingAccount: "billingAccounts/54321",
			}
		},
	})
}

func Test_FreightService_AlternativeSetup(t *testing.T) {
	// Even though no implementation exists, the tests will pass but be skipped.
	freightv1aiptest.TestServices(t, &aipTests{})
}

type aipTests struct{}

var _ freightv1aiptest.FreightServiceTestSuiteConfigProvider = &aipTests{}

func (a aipTests) FreightServiceShipper(_ *testing.T) *freightv1aiptest.FreightServiceShipperTestSuiteConfig {
	// Returns nil to indicate that it's not ready to be tested.
	return nil
}

func (a aipTests) FreightServiceSite(_ *testing.T) *freightv1aiptest.FreightServiceSiteTestSuiteConfig {
	// Returns nil to indicate that it's not ready to be tested.
	return nil
}
