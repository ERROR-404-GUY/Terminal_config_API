package domain_test

import (
	"context"
	"terminal_config/internal/domain"
	"terminal_config/internal/domain/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestCreateTerminal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTerminalRepository(ctrl)

	service := domain.NewTerminalConfigService(mockRepo)

	// table driven tests
	// allow you to test for loads of variations of input easily
	tests := []struct {
		name    string
		input   *domain.Terminal
		mocking func()
		wantErr bool
	}{
		{
			name:  "Success case",
			input: &domain.Terminal{SerialNumber: "SN-123"},
			mocking: func() {
				mockRepo.EXPECT().
					CreateTerminal(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		// add more test cases here!

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocking()
			err := service.CreateTerminal(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error = %v, but got %v", tt.wantErr, err)
			}
		})
	}
}
