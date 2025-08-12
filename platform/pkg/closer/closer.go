package closer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
)

const shutdownTimeout = 5 * time.Second

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Closer struct {
	mu     sync.Mutex
	once   sync.Once
	done   chan struct{}
	funcs  []func(ctx context.Context) error
	logger Logger
}

// Глобальный экземпляр для использования по всему приложению
var globalCloser = NewWithLogger(&logger.NoopLogger{})

// AddNamed добавляет функцию закрытия с именем зависимости для логирования в глобальный closer
func AddNamed(name string, f func(context.Context) error) {
	globalCloser.AddNamed(name, f)
}

// Add добавляет функции закрытия в глобальный closer
func Add(f ...func(context.Context) error) {
	globalCloser.Add(f...)
}

// CloseAll инициирует процесс закрытия всех зарегистрированных функций глобального closer'а
func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

// SetLogger позволяет установить кастомный логгер для глобального closer'а
func SetLogger(l Logger) {
	globalCloser.SetLogger(l)
}

// Configure настраивает глобальный closer для обработки системных сигналов
func Configure(signals ...os.Signal) {
	go globalCloser.handleSignals(signals...)
}

// New создаёт новый экземпляр Closer с дефолтным логгером log.Default()
func New(signals ...os.Signal) *Closer {
	return NewWithLogger(logger.Logger(), signals...)
}

func NewWithLogger(logger Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done:   make(chan struct{}),
		logger: logger,
	}

	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}

	return c
}

func (c *Closer) SetLogger(l Logger) {
	c.logger = l
}

func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		c.logger.Info(context.Background(), "🛑 Received system signal, initializing graceful shutdown")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.logger.Error(context.Background(), "❌ Error has occurred while closing resources: %v", zap.Error(err))
		}

	case <-c.done:

	}
}

func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()
		c.logger.Info(ctx, fmt.Sprintf("🧩 Closing %s...", name))

		err := f(ctx)

		duration := time.Since(start)
		if err != nil {
			c.logger.Error(ctx, fmt.Sprintf("❌ Error while closing %s: %v (took %s)", name, err, duration))
		} else {
			c.logger.Info(ctx, fmt.Sprintf("✅ %s closed in %s", name, duration))
		}
		return err
	})
}

func (c *Closer) Add(f ...func(ctx context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.logger.Info(ctx, "ℹ️ No funcs to close")
			return
		}

		c.logger.Info(ctx, "🚦 starting graceful shutdown")

		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup

		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]
			wg.Add(1)
			go func(f func(ctx2 context.Context) error) {
				defer wg.Done()

				defer func() {
					if r := recover(); r != nil {
						errCh <- errors.New("panic recovered in closer")
						c.logger.Error(ctx, "⚠️ Panic in closer func", zap.Any("error", r))
					}
				}()

				if err := f(ctx); err != nil {
					errCh <- err
				}
			}(f)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for {
			select {
			case <-ctx.Done():
				c.logger.Info(ctx, "⚠️ context cancelled while closing", zap.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}
				return
			case err, ok := <-errCh:
				if !ok {
					c.logger.Info(ctx, "✅ All resources successfully closed")
					return
				}
				c.logger.Error(ctx, "❌ Error while closing", zap.Error(err))
				if result == nil {
					result = err
				}
			}
		}
	})

	return result
}
