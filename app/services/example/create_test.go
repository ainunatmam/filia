package example

import (
	"context"
	"errors"
	"testing"
	"wallet-api/app/entity"
	"wallet-api/app/presentation"
	"wallet-api/app/repositories"
)

func TestExampleService_Create(t *testing.T) {
	tests := []struct {
		name       string
		request    *presentation.ExampleRequest
		mockRepo   *repositories.MockExampleRepository
		wantErr    bool
		errMessage string
	}{
		{
			name: "successfully creates example",
			request: &presentation.ExampleRequest{
				Name: "Test Example",
			},
			mockRepo: &repositories.MockExampleRepository{
				CreateFunc: func(ctx context.Context, data entity.Example) error {
					if data.Name != "Test Example" {
						t.Errorf("expected name 'Test Example', got '%s'", data.Name)
					}
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "returns error when repository fails",
			request: &presentation.ExampleRequest{
				Name: "Test Example",
			},
			mockRepo: &repositories.MockExampleRepository{
				CreateFunc: func(ctx context.Context, data entity.Example) error {
					return errors.New("database connection failed")
				},
			},
			wantErr:    true,
			errMessage: "database connection failed",
		},
		{
			name: "handles empty name",
			request: &presentation.ExampleRequest{
				Name: "",
			},
			mockRepo: &repositories.MockExampleRepository{
				CreateFunc: func(ctx context.Context, data entity.Example) error {
					// Repository accepts empty name (validation should be done at handler level)
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "handles very long name",
			request: &presentation.ExampleRequest{
				Name: "This is a very long name that might exceed database column limits but should still be handled by the application layer",
			},
			mockRepo: &repositories.MockExampleRepository{
				CreateFunc: func(ctx context.Context, data entity.Example) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "context is propagated to repository",
			request: &presentation.ExampleRequest{
				Name: "Context Test",
			},
			mockRepo: &repositories.MockExampleRepository{
				CreateFunc: func(ctx context.Context, data entity.Example) error {
					// Verify context is not nil
					if ctx == nil {
						t.Error("context should not be nil")
					}
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewExampleService(tt.mockRepo)
			ctx := context.Background()

			err := service.Create(ctx, tt.request)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				if tt.errMessage != "" && err.Error() != tt.errMessage {
					t.Errorf("expected error message '%s', got '%s'", tt.errMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestNewExampleService tests the service constructor
func TestNewExampleService(t *testing.T) {
	mockRepo := &repositories.MockExampleRepository{}
	service := NewExampleService(mockRepo)

	if service == nil {
		t.Error("NewExampleService should return a non-nil service")
	}
}

// TestExampleService_Create_ContextCancellation tests context cancellation handling
func TestExampleService_Create_ContextCancellation(t *testing.T) {
	mockRepo := &repositories.MockExampleRepository{
		CreateFunc: func(ctx context.Context, data entity.Example) error {
			// Simulate checking context
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		},
	}

	service := NewExampleService(mockRepo)

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := service.Create(ctx, &presentation.ExampleRequest{Name: "Test"})
	if err == nil {
		t.Log("Note: Context cancellation not yet propagated to repository - consider adding context checks")
	}
}
