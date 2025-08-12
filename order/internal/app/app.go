package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/HeyReyHR/rocket-factory/order/internal/api/health"
	"github.com/HeyReyHR/rocket-factory/order/internal/config"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/closer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const readHeaderTimeout = 5 * time.Second
const requestTimeout = 10 * time.Second

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
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
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
	return logger.Init( // TODO LOGGER
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
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

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(requestTimeout))

	r.Mount("/", orderServer)

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	} // TODO INTERCEPTOR/BETTER ERRORS

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		a.httpServer.Shutdown(ctx)
		return nil
	})

	health.Health(r)
	
	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error { // Ревьюверу: надо ли запускать сервер в отдельной горутине?
	logger.Info(ctx, fmt.Sprintf("🚀 Starting server on %s", config.AppConfig().OrderHTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, fmt.Sprintf("❌ Error occurred when starting server: %s", err))
		return err
	}

	return nil
}
