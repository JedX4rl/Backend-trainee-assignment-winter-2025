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
	Id        int    `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Coins     int64  `json:"coins" db:"coins"`
	Inventory Inventory
}
