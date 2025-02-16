package models

type CoinHistory struct {
	Received []TransactionInfo `json:"received"`
	Sent     []TransactionInfo `json:"sent"`
}

type InfoResponse struct {
	Coins       int64           `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coin_history"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

func (e ErrorResponse) Error() string {
	return e.Errors
}

type AuthResponse struct {
	Token string `json:"token"`
}
