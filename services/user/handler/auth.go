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
		token, err := auth.GenerateAccessToken(u.ID, u.IsAdmin, ttl)
		if err != nil {
			h.log.Error("failed to generate token", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(loginResp{AccessToken: token, ExpiresIn: int64(ttl.Seconds())})
	}
}
