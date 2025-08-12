package e2e

import (
	"context"
	"fmt"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const testsTimeout = 5 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Payment Service Integration Test Suite")
}

var _ = BeforeSuite(func() {

	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("Cannot init logger: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	envVars, err := godotenv.Read(filepath.Join("..", "..", "..", "deploy", "env", ".env"))
	if err != nil {
		logger.Fatal(suiteCtx, "Cannot load .env", zap.Error(err))
	}

	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

	logger.Info(suiteCtx, "Starting test environment...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info(context.Background(), "Shutting down tests")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}

	suiteCancel()
})
