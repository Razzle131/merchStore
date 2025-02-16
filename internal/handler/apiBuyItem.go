package handler

import (
	"log/slog"
	"net/http"

	"github.com/Razzle131/merchStore/internal/serverErrors"
)

// Купить предмет за монеты.
// (GET /api/buy/{item})
func (s *MyServer) GetApiBuyItem(w http.ResponseWriter, r *http.Request, item string) {
	slog.Info("proccessing buy item")
	defer slog.Info("finished buy item")

	token := r.Header.Get("Authorization")

	user, err := s.auth.AuthorizeUser(token)
	if err != nil {
		sendErrorResponse(w, "authorization process failed", http.StatusUnauthorized)
		return
	}

	err = s.payment.BuyMerch(user.Id, item)
	if err != nil && (err == serverErrors.ErrNotAllowed || err == serverErrors.ErrUserNotFound || err == serverErrors.ErrItemNotFound) {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		sendErrorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug("purchase is successful")

	sendInfoResponse(w, nil)
}
