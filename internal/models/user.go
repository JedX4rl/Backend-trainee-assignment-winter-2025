package models

type Inventory struct {
	Type     string
	Quantity int64
}

//TODO rename function

type transactionInfo struct {
	Username string
	Amount   int32
}

type CoinHistory struct {
	Received []transactionInfo
	Sent     []transactionInfo
}

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Coins     int64  `json:"coins"`
	Inventory Inventory
}
