package models

import "encoding/json"

type TransactionInfo struct {
	Username string `json:"-"`
	Amount   int32  `json:"amount"`
	IsSent   bool   `json:"-"`
}

func (t TransactionInfo) MarshalJSON() ([]byte, error) {
	alias := struct {
		FromUser string `json:"fromUser,omitempty"`
		ToUser   string `json:"toUser,omitempty"`
		Amount   int32  `json:"amount"`
	}{
		Amount: t.Amount,
	}

	if t.IsSent {
		alias.ToUser = t.Username
	} else {
		alias.FromUser = t.Username
	}

	return json.Marshal(alias)
}
