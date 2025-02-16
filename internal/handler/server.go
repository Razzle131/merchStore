package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Razzle131/merchStore/api"
	"github.com/Razzle131/merchStore/internal/db"
	"github.com/Razzle131/merchStore/internal/repository/merchRepo"
	"github.com/Razzle131/merchStore/internal/repository/transactionRepo"
	"github.com/Razzle131/merchStore/internal/repository/userRepo"
	"github.com/Razzle131/merchStore/internal/service/authorization"
	"github.com/Razzle131/merchStore/internal/service/payment"
)

type MyServer struct {
	auth    authorization.AuthorizationService
	payment payment.PaymentService
}

type Config struct {
	Port string
	DSN  string
}

var _ api.ServerInterface = (*MyServer)(nil)

func NewServer(db *db.DB) *MyServer {
	userRepo := userRepo.NewUserRepoPg(db)
	merchRepo := merchRepo.NewMerchRepoPg(db)
	transactionRepo := transactionRepo.NewTransactionRepoPg(db)

	return &MyServer{
		auth:    *authorization.New(userRepo),
		payment: *payment.New(userRepo, merchRepo, transactionRepo),
	}
}

func sendErrorResponse(w http.ResponseWriter, errMsg string, status int) {
	resp, _ := json.Marshal(api.ErrorResponse{Errors: &errMsg})
	slog.Error(errMsg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}

func sendInfoResponse(w http.ResponseWriter, object any) {
	if object != nil {
		resp, err := json.Marshal(object)
		if err != nil {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// to ensure that object is converted and there are no error and we dont have "superfluous response.WriteHeader call" message in log
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
}
