package model

type WalletInfo struct {
	In   []Transaction
	Out  []Transaction
	Cash int
}

func NewWallet() WalletInfo {
	return WalletInfo{
		In:   make([]Transaction, 0, 16),
		Out:  make([]Transaction, 0, 16),
		Cash: 1000,
	}
}
