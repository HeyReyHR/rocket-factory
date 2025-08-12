package health

import (
	"net/http"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func Health(r chi.Router) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.Context(), "💓 health check called")

		w.WriteHeader(http.StatusOK)
		// nolint:errcheck,gosec
		w.Write([]byte("ok"))
	})
}
