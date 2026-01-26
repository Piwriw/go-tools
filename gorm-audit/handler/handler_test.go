package handler

import (
	"context"
	"testing"
)

func TestEventHandlerFunc(t *testing.T) {
	called := false
	fn := EventHandlerFunc(func(ctx context.Context, event *Event) error {
		called = true
		return nil
	})

	err := fn.Handle(context.Background(), &Event{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("handler was not called")
	}
}

func TestEventHandlerFuncWithError(t *testing.T) {
	expectedErr := &TestError{}
	fn := EventHandlerFunc(func(ctx context.Context, event *Event) error {
		return expectedErr
	})

	err := fn.Handle(context.Background(), &Event{})
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

type TestError struct{}

func (e *TestError) Error() string {
	return "test error"
}
