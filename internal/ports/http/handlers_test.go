package ports_test

import (
	"net/http"
	"net/http/httptest"
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
}
