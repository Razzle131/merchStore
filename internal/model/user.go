package model

type User struct {
	Id       int
	Login    string
	Password string
	Wallet   WalletInfo
	Items    Inventory
}

func NewUser(id int, login, password string) User {
	return User{
		Id:       id,
		Login:    login,
		Password: password,
		Wallet:   NewWallet(),
		Items:    NewInventory(),
	}
}
