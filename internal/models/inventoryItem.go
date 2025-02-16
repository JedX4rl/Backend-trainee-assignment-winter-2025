package models

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int64  `json:"quantity"`
}
