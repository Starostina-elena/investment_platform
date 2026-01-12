package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler := AuthMiddleware(nextHandler)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
	})

	tests := []struct {
		name   string
		header string
	}{
		{"No Bearer", "token123"},
		{"Wrong prefix", "Basic token123"},
		{"Empty token", "Bearer "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.header)
			w := httptest.NewRecorder()

			handler := AuthMiddleware(nextHandler)
			handler.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected status 401, got %d", w.Code)
			}
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()

	handler := AuthMiddleware(nextHandler)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestFromContext_Found(t *testing.T) {
	expected := &UserClaims{
		UserID: 123,
		Admin:  true,
		Banned: false,
	}

	ctx := context.WithValue(context.Background(), ctxUserKey, expected)

	result := FromContext(ctx)
	if result == nil {
		t.Fatal("expected UserClaims, got nil")
	}

	if result.UserID != expected.UserID {
		t.Errorf("expected UserID %d, got %d", expected.UserID, result.UserID)
	}

	if result.Admin != expected.Admin {
		t.Errorf("expected Admin %v, got %v", expected.Admin, result.Admin)
	}
}

func TestFromContext_NotFound(t *testing.T) {
	ctx := context.Background()

	result := FromContext(ctx)
	if result != nil {
		t.Errorf("expected nil, got %+v", result)
	}
}

func TestFromContext_WrongType(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxUserKey, "wrong type")

	result := FromContext(ctx)
	if result != nil {
		t.Errorf("expected nil, got %+v", result)
	}
}
