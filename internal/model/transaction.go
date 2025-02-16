package model

type Transaction struct {
	LoginTo   string
	LoginFrom string
	Amount    int
}

func NewTransaction(to, from string, amount int) Transaction {
	return Transaction{
		LoginTo:   to,
		LoginFrom: from,
		Amount:    amount,
	}
}
