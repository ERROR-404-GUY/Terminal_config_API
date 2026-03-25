package ports

import (
	"encoding/json"
	"net/http"
	"terminal_config/internal/domain"

	"strconv"

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
	h.muxRouter.HandleFunc("/api/terminals", h.CreateTerminalHandler).Methods("POST")
	h.muxRouter.HandleFunc("/api/terminals/random", h.CreateRandomTerminalHandler).Methods("POST")
	h.muxRouter.HandleFunc("/api/terminals/{tid}", h.GetTerminalHandler).Methods("GET")
	h.muxRouter.HandleFunc("/api/terminals/{tid}", h.UpdateTerminalHandler).Methods("PUT")
	h.muxRouter.HandleFunc("/api/terminals/{tid}", h.DeleteTerminalHandler).Methods("DELETE")
	h.muxRouter.HandleFunc("/api/terminals", h.ListTerminalsHandler).Methods("GET")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (h *handlers) CreateTerminalHandler(w http.ResponseWriter, r *http.Request) {
	var terminal domain.Terminal
	if err := json.NewDecoder(r.Body).Decode(&terminal); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Check for refund_allowed query parameter
	refundAllowedStr := r.URL.Query().Get("refund_allowed")
	if refundAllowedStr != "" {
		refundAllowed, err := strconv.ParseBool(refundAllowedStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid refund_allowed value"})
			return
		}
		terminal.RefundAllowed = refundAllowed
	}

	if err := h.service.CreateTerminal(r.Context(), &terminal); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(terminal)
}

func (h *handlers) GetTerminalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tid := vars["tid"]

	terminal, err := h.service.GetTerminal(r.Context(), tid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Terminal not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(terminal)
}

func (h *handlers) UpdateTerminalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tid := vars["tid"]

	// Check for refund_allowed query parameter
	refundAllowedStr := r.URL.Query().Get("refund_allowed")
	var refundAllowed *bool
	if refundAllowedStr != "" {
		parsed, err := strconv.ParseBool(refundAllowedStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid refund_allowed value"})
			return
		}
		refundAllowed = &parsed
	}

	// If no body and only refund_allowed query param, do targeted update
	if r.Body == nil || r.ContentLength == 0 {
		if refundAllowed != nil {
			if err := h.service.UpdateRefundAllowed(r.Context(), tid, *refundAllowed); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			// Return the updated terminal
			terminal, err := h.service.GetTerminal(r.Context(), tid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(terminal)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "No body or query parameters provided"})
			return
		}
	}

	// Body provided - decode and update
	var terminal domain.Terminal
	if err := json.NewDecoder(r.Body).Decode(&terminal); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	terminal.TID = tid

	// Override refund_allowed if query param provided
	if refundAllowed != nil {
		terminal.RefundAllowed = *refundAllowed
	}

	if err := h.service.UpdateTerminal(r.Context(), &terminal); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(terminal)
}

func (h *handlers) DeleteTerminalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tid := vars["tid"]

	if err := h.service.DeleteTerminal(r.Context(), tid); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) CreateRandomTerminalHandler(w http.ResponseWriter, r *http.Request) {
	terminal, err := h.service.CreateRandomTerminal(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(terminal)
}

func (h *handlers) ListTerminalsHandler(w http.ResponseWriter, r *http.Request) {
	terminals, err := h.service.ListTerminals(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(terminals)
}
