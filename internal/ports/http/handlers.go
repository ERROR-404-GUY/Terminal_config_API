package ports

import (
	"net/http"
	"terminal_config/internal/domain"

	"github.com/gorilla/mux"
)

type Handlers interface {
	SetupRoutes()
}
type handlers struct {
	service   domain.TerminalConfigService
	muxRouter *mux.Router
}

func NewHandlers(muxRouter *mux.Router, service domain.TerminalConfigService) Handlers {
	h := &handlers{
		service:   service,
		muxRouter: muxRouter,
	}

	// Set up your HTTP routes here using muxRouter and h's methods

	return h
}

func (h *handlers) SetupRoutes() {
	h.muxRouter.HandleFunc("/api/health", HealthCheckHandler).Methods("GET")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
