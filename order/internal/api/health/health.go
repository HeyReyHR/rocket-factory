package health

import (
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Health(r chi.Router) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.Context(), "💓 health check called")
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}
