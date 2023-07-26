package models

type TransactionRequestPayload struct {
	UID    string `json:"UID"`
	ItemId string `json:"itemid"`
	Count  int    `json:"count"`
}
type Transaction struct {
	UID             string `json:"UID"`
	ItemId          string `json:"itemid"`
	TransactionId   string `json:"transactionId"`
	TransactionDate string `json:"transactionDate"`
	Count           string `json:"count"`
	IsFinished      bool   `json:"isFinished"`
}
