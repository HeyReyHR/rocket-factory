package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/HeyReyHR/rocket-factory/order/internal/api/health"
	"github.com/HeyReyHR/rocket-factory/order/internal/config"
	orderMetrics "github.com/HeyReyHR/rocket-factory/order/internal/metrics"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/closer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/metrics"
	auth "github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/http"
	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

const (
	readHeaderTimeout = 5 * time.Second
	requestTimeout    = 10 * time.Second
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("consumer crashed: %v", err)
		}
	}()

	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- fmt.Errorf("http server crashed: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		cancel()
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initMetricsProvider,
		a.initMetrics,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
		config.AppConfig().Logger.EnableOTLP(),
		config.AppConfig().Logger.OTLPEnvironment(),
		config.AppConfig().Logger.OTLPServiceName(),
	)
}

func (a *App) initMetricsProvider(ctx context.Context) error {
	return metrics.InitProvider(ctx, config.AppConfig().Metrics)
}

func (a *App) initMetrics(_ context.Context) error {
	return orderMetrics.InitMetrics()
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	orderAPI := a.diContainer.OrderV1API(ctx)

	orderServer, err := orderV1.NewServer(orderAPI)
	if err != nil {
		logger.Error(ctx, "❌ Error occurred when creating OpenAPI server", zap.Error(err))
	}

	r := chi.NewRouter()

	r.Use(auth.NewAuthMiddleware(a.diContainer.IamClient(ctx)).Handle)
	r.Use(middleware.Recoverer, middleware.Logger)
	r.Use(middleware.Timeout(requestTimeout))

	r.Mount("/", orderServer)

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		err = a.httpServer.Shutdown(ctx)
		if err != nil {
			logger.Error(ctx, "Cannot shutdown server")
		}
		return nil
	})

	health.Health(r)

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 Starting server on %s", config.AppConfig().OrderHTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, fmt.Sprintf("❌ Error occurred when starting server: %s", err))
		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 ShipAssembled Kafka consumer running")

	err := a.diContainer.ShipConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
