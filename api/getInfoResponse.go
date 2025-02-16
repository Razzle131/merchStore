package api

type GetInfoRecieved struct {
	Amount   *int    `json:"amount,omitempty"`
	FromUser *string `json:"fromUser,omitempty"`
}

type GetInfoSent struct {
	Amount *int    `json:"amount,omitempty"`
	ToUser *string `json:"toUser,omitempty"`
}

type GetInfoCoinHistory struct {
	Received *[]GetInfoRecieved `json:"received,omitempty"`
	Sent     *[]GetInfoSent     `json:"sent,omitempty"`
}

type GetInfoInventory struct {
	Quantity *int    `json:"quantity,omitempty"`
	Type     *string `json:"type,omitempty"`
}

type GetInfoResponse struct {
	CoinHistory *GetInfoCoinHistory `json:"coinHistory,omitempty"`
	Coins       *int                `json:"coins,omitempty"`
	Inventory   *[]GetInfoInventory `json:"inventory,omitempty"`
}
