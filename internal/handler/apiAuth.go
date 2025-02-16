package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Razzle131/merchStore/api"
	"github.com/Razzle131/merchStore/internal/serverErrors"
)

// Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
// (POST /api/auth)
func (s *MyServer) PostApiAuth(w http.ResponseWriter, r *http.Request) {
	slog.Info("proccessing auth")
	defer slog.Info("finished auth")

	var authStruct api.AuthRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&authStruct)
	if err != nil {
		sendErrorResponse(w, "bad request body", http.StatusBadRequest)
		return
	}
	slog.Debug("body decoded")

	token, err := s.auth.AuthenticateUser(authStruct.Username, authStruct.Password)
	if err != nil && err == serverErrors.ErrBadCreditonals {
		sendErrorResponse(w, "bad creditionals", http.StatusBadRequest)
		return
	} else if err != nil {
		sendErrorResponse(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := api.AuthResponse{
		Token: &token,
	}

	sendInfoResponse(w, resp)
}
