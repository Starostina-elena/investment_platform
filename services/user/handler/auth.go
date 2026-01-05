package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Starostina-elena/investment_platform/services/user/auth"
	"github.com/Starostina-elena/investment_platform/services/user/core"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func LoginHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		u, err := h.service.GetByEmail(r.Context(), req.Email)
		if err != nil {
			h.log.Error("login: user not found", "email", req.Email, "error", err)
			http.Error(w, "Некорретный email", http.StatusUnauthorized)
			return
		}

		if err := core.VerifyPassword(u.PasswordHash, req.Password); err != nil {
			h.log.Info("login: bad password", "user", u.ID)
			http.Error(w, "Некорретный пароль", http.StatusUnauthorized)
			return
		}

		ttl := 15 * time.Minute
		token, err := auth.GenerateAccessToken(u.ID, u.IsAdmin, u.IsBanned, ttl)
		if err != nil {
			h.log.Error("failed to generate token", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		refreshTTL := 7 * 24 * time.Hour
		rawRefresh, err := h.service.GenerateRefreshToken(r.Context(), u.ID, refreshTTL)
		if err != nil {
			h.log.Error("failed to create refresh token", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    rawRefresh,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(refreshTTL),
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(loginResp{AccessToken: token, ExpiresIn: int64(ttl.Seconds())})
	}
}

type refreshResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func RefreshHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		raw := cookie.Value

		u, err := h.service.AuthenticateByRefresh(r.Context(), raw)
		if err != nil {
			h.log.Info("refresh failed", "error", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ttl := 15 * time.Minute
		accessToken, err := auth.GenerateAccessToken(u.ID, u.IsAdmin, u.IsBanned, ttl)
		if err != nil {
			h.log.Error("failed to generate access token", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(refreshResp{AccessToken: accessToken, ExpiresIn: int64(ttl.Seconds())})
	}
}

func LogoutHandler(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err == nil {
			hash := auth.HashToken(cookie.Value)
			_ = h.service.RevokeRefreshToken(r.Context(), hash)
			clear := &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
			}
			http.SetCookie(w, clear)
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
