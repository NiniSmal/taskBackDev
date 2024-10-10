package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"taskBackDev/config"
	"taskBackDev/entity"
	"taskBackDev/service"
	"taskBackDev/service/send_service"
	"time"
)

type Handler struct {
	service *service.Service
	cfg     *config.Config
	ss      *send_service.SendService
}

func NewHandler(s *service.Service, cfg *config.Config, ss *send_service.SendService) *Handler {
	return &Handler{service: s, cfg: cfg, ss: ss}
}

const (
	name          = "Name"
	authorisation = "Authorization"
	email         = "email"
)

func (h *Handler) AuthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idP := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(idP, 10, 64)
		if err != nil {
			http.Error(w, "parsing to id", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userName, err := getHeader(r, name)
		if err != nil {
			http.Error(w, fmt.Sprintf("Header '%v' is missing", name), http.StatusBadRequest)
			return
		}

		token, err := h.service.Tokens(ctx, userName, id)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusUnauthorized)
			return
		}

		refreshToken := strconv.Itoa(int(token.RefreshToken))
		setCookies(w, refreshToken, h.cfg.RefreshTokenTTL)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(token)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusBadRequest)

			return
		}
	}
}
func getHeader(r *http.Request, header string) (string, error) {
	h := r.Header.Get(header)
	if h == "" {
		return "", fmt.Errorf(fmt.Sprintf("Header '%v' is missing", header))
	}

	return h, nil
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	header := r.Header.Get(authorisation)
	if header == "" {
		w.Write([]byte("empty auth header"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		http.Error(w, "invalid auth header", http.StatusBadRequest)
		return
	}

	if len(headerParts[1]) == 0 {
		http.Error(w, "ip is empty", http.StatusBadRequest)
		return
	}
	ipName, err := h.service.Parse(headerParts[1])
	if err != nil {
		http.Error(w, "parse access token by ip", http.StatusBadRequest)
		return
	}
	userName, err := getHeader(r, name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Header '%v' is missing", name), http.StatusBadRequest)
		return
	}
	if ipName != userName {
		email := entity.Email{
			Text:    "Warning",
			To:      email,
			Subject: "Отчет за месяц",
		}
		err = h.ss.SendMessage(email)
		if err != nil {
			http.Error(w, "send message", http.StatusInternalServerError)
			return
		}
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "get cookie", http.StatusUnauthorized)
		return
	}

	userID, err := h.service.SwitchToken(ctx, cookie.Value, ipName)
	if err != nil {
		http.Error(w, "switch token", http.StatusUnauthorized)
		return
	}

	res, err := h.service.CreateToken(ctx, userID, ipName)
	if err != nil {
		if errors.Is(err, entity.ErrUserBlocked) {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(err.Error()))
		return
	}
	refreshToken := strconv.Itoa(int(res.RefreshToken))
	setCookies(w, refreshToken, h.cfg.RefreshTokenTTL)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}

func setCookies(w http.ResponseWriter, refreshToken string, refreshTokenTTL time.Duration) {
	httpOnlyCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(refreshTokenTTL),
		Path:     "/api/auth",
		HttpOnly: true,
	}
	http.SetCookie(w, &httpOnlyCookie)
}
