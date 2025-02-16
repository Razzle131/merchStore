package handler

import (
	"log/slog"
	"net/http"

	"github.com/Razzle131/merchStore/api"
	"github.com/Razzle131/merchStore/internal/serverErrors"
)

// Получить информацию о монетах, инвентаре и истории транзакций.
// (GET /api/info)
func (s *MyServer) GetApiInfo(w http.ResponseWriter, r *http.Request) {
	slog.Info("proccessing info")
	defer slog.Info("finished info")

	token := r.Header.Get("Authorization")
	user, err := s.auth.AuthorizeUser(token)
	if err != nil {
		sendErrorResponse(w, "authorization process failed", http.StatusUnauthorized)
		return
	}

	resCoinHistory := api.GetInfoCoinHistory{}
	resRecieved := make([]api.GetInfoRecieved, 0, 16)
	resSent := make([]api.GetInfoSent, 0, 16)
	resInventory := make([]api.GetInfoInventory, 0, 16)

	wallet, err := s.payment.GetWalletInfo(user.Id)
	if err != nil && err == serverErrors.ErrUserNotFound {
		sendErrorResponse(w, "no such user", http.StatusBadRequest)
		return
	} else if err != nil {
		sendErrorResponse(w, "server internal error", http.StatusInternalServerError)
		return
	}
	slog.Debug("got wallet")

	for _, inTransaction := range wallet.In {
		resRecieved = append(resRecieved, api.GetInfoRecieved{
			Amount:   &inTransaction.Amount,
			FromUser: &inTransaction.LoginFrom,
		})
	}

	for _, outTransaction := range wallet.Out {
		resSent = append(resSent, api.GetInfoSent{
			Amount: &outTransaction.Amount,
			ToUser: &outTransaction.LoginTo,
		})
	}

	for item, quantity := range user.Items.GetItems() {
		resInventory = append(resInventory, api.GetInfoInventory{
			Quantity: &quantity,
			Type:     &item,
		})
	}

	resCoinHistory.Received = &resRecieved
	resCoinHistory.Sent = &resSent

	res := api.GetInfoResponse{}
	res.CoinHistory = &resCoinHistory
	res.Coins = &user.Wallet.Cash
	res.Inventory = &resInventory

	sendInfoResponse(w, res)
}
