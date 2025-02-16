package transactionRepo

import "github.com/Razzle131/merchStore/internal/model"

type TransactionRepoInterface interface {
	// mode 1 -> income transactions
	//
	// mode 2 -> sent transactions
	GetTransactions(login string, mode int) ([]model.Transaction, error)

	AddTransaction(loginFrom, loginTo string, amount int) error
}
