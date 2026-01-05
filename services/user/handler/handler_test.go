package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Starostina-elena/investment_platform/services/user/core"
)

func TestCreateUserHandler_Success(t *testing.T) {
	h := setupTestHandler(t)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"nickname": "testuser",
		"password": "password123",
		"name":     "Test",
		"surname":  "User",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/create", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler := CreateUserHandler(h)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateUserHandler_InvalidJSON(t *testing.T) {
	h := setupTestHandler(t)

	req := httptest.NewRequest("POST", "/create", bytes.NewBuffer([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler := CreateUserHandler(h)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGetUserHandler_InvalidID(t *testing.T) {
	h := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler := GetUserHandler(h)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateUserHandler_Unauthorized(t *testing.T) {
	h := setupTestHandler(t)

	reqBody := map[string]string{
		"email": "updated@example.com",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/update", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler := UpdateUserHandler(h)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func setupTestHandler(_ *testing.T) *Handler {
	return &Handler{
		service: &mockService{},
		log:     *mockLogger(),
	}
}

type mockService struct{}

func (m *mockService) GetByEmail(ctx context.Context, email string) (*core.User, error) {
	return nil, core.ErrUserNotFound
}

func (m *mockService) GenerateRefreshToken(ctx context.Context, userID int, ttl time.Duration) (string, error) {
	return "", nil
}

func (m *mockService) AuthenticateByRefresh(ctx context.Context, token string) (*core.User, error) {
	return nil, core.ErrInvalidToken
}

func (m *mockService) RevokeRefreshToken(ctx context.Context, hash string) error {
	return nil
}

func (m *mockService) Create(ctx context.Context, user core.User) (*core.User, error) {
	return &core.User{ID: 1, Email: user.Email, Nickname: user.Nickname}, nil
}

func (m *mockService) Get(ctx context.Context, id int) (*core.User, error) {
	return nil, core.ErrUserNotFound
}

func (m *mockService) Update(ctx context.Context, user core.User) (*core.User, error) {
	return &user, nil
}

func (m *mockService) SetAdmin(ctx context.Context, userID int, isAdmin bool) error {
	return nil
}

func (m *mockService) BanUser(ctx context.Context, userID int, isBanned bool) error {
	return nil
}

func mockLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
}
