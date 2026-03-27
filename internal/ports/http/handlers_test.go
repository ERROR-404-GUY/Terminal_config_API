package ports_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"terminal_config/internal/domain"
	"terminal_config/internal/domain/mocks"
	ports "terminal_config/internal/ports/http"
	"testing"

	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

func TestGetTerminalHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockTerminalConfigService(ctrl)
	router := mux.NewRouter()

	h := ports.NewHandlers(router, mockSvc)
	h.SetupRoutes()

	t.Run("Returns 200 when terminal exists", func(t *testing.T) {
		mockSvc.EXPECT().
			GetTerminal(gomock.Any(), "TID-123").
			Return(nil, nil)

		req := httptest.NewRequest("GET", "/api/terminals/TID-123", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusOK)
		}
	})
	t.Run("Returns 404 when terminal does not exist", func(t *testing.T) {
		mockSvc.EXPECT().
			GetTerminal(gomock.Any(), "TID-134").
			Return(nil, errors.New("not found"))

		req := httptest.NewRequest("GET", "/api/terminals/TID-134", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusNotFound)
		}
	})

}

func TestCreateTerminalHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockTerminalConfigService(ctrl)
	router := mux.NewRouter()

	h := ports.NewHandlers(router, mockSvc)
	h.SetupRoutes()

	t.Run("Returns 201 when created successfully", func(t *testing.T) {
		mockSvc.EXPECT().
			CreateTerminal(gomock.Any(), gomock.Any()).
			Return(nil)

		payload := `{"serial_number": "SN-CREATE-123", "product_name": "Test Terminal"}`
		req := httptest.NewRequest("POST", "/api/terminals", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusCreated)
		}
	})

	t.Run("Returns 400 when invalid body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/terminals", bytes.NewBufferString(`{invalid json`))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusBadRequest)
		}
	})
}

func TestDeleteTerminalHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockTerminalConfigService(ctrl)
	router := mux.NewRouter()

	h := ports.NewHandlers(router, mockSvc)
	h.SetupRoutes()

	t.Run("Returns 204 when deleted successfully", func(t *testing.T) {
		mockSvc.EXPECT().
			DeleteTerminal(gomock.Any(), "TID-DELETE-123").
			Return(nil)

		req := httptest.NewRequest("DELETE", "/api/terminals/TID-DELETE-123", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusNoContent)
		}
	})
}

func TestListTerminalsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockTerminalConfigService(ctrl)
	router := mux.NewRouter()

	h := ports.NewHandlers(router, mockSvc)
	h.SetupRoutes()

	t.Run("Returns 200 and a list of terminals", func(t *testing.T) {
		mockSvc.EXPECT().
			ListTerminals(gomock.Any()).
			Return([]*domain.Terminal{{TID: "TID-LIST-1"}}, nil)

		req := httptest.NewRequest("GET", "/api/terminals", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusOK)
		}
	})
}

func TestUpdateTerminalHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockTerminalConfigService(ctrl)
	router := mux.NewRouter()

	h := ports.NewHandlers(router, mockSvc)
	h.SetupRoutes()

	t.Run("Returns 200 if update is successful", func(t *testing.T) {
		mockSvc.EXPECT().
			UpdateRefundAllowed(gomock.Any(), "TID-UPDATE-123", true).
			Return(nil)

		mockSvc.EXPECT().
			GetTerminal(gomock.Any(), "TID-UPDATE-123").
			Return(&domain.Terminal{TID: "TID-UPDATE-123", RefundAllowed: true}, nil)

		req := httptest.NewRequest("PUT", "/api/terminals/TID-UPDATE-123?refund_allowed=true", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusOK)
		}
	})
}
