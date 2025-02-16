package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Razzle131/merchStore/api"
	"github.com/Razzle131/merchStore/internal/serverErrors"
)

// Отправить монеты другому пользователю.
// (POST /api/sendCoin)
func (s *MyServer) PostApiSendCoin(w http.ResponseWriter, r *http.Request) {
	slog.Info("proccessing send coin")
	defer slog.Info("finished send coin")

	token := r.Header.Get("Authorization")
	user, err := s.auth.AuthorizeUser(token)
	if err != nil {
		sendErrorResponse(w, "authorization process failed", http.StatusUnauthorized)
		return
	}

	var req api.SendCoinRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&req)
	if err != nil {
		sendErrorResponse(w, "bad request", http.StatusBadRequest)
		return
	}
	slog.Debug("body decoded")

	if req.ToUser == user.Login {
		sendErrorResponse(w, "cant send coins to yourself", http.StatusBadRequest)
		return
	}
	err = s.payment.SendCoin(req.ToUser, user.Login, req.Amount)
	if err != nil && err == serverErrors.ErrUserNotFound {
		sendErrorResponse(w, "reciever or sender not found", http.StatusBadRequest)
		return
	} else if err != nil {
		sendErrorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	sendInfoResponse(w, nil)
}
