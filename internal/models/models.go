package models

import "time"

type Item struct {
	ID          string `json:"id"`
	Price       int    `json:"price"`
	Name        string `json:"Name"`
	Stonesize   int    `json:"Stonesize"`
	Minlength   int    `json:"Minlength"`
	Maxlength   int    `json:"Maxlength"`
	Decsription string `json:"Description"`
	Keyword     string `json:"Keyword"`
}

type CustomerRequestPayload struct {
	UID         string `json:"uid"`
	CreatedDate time.Time
	Email       string `json:"contact"`
}
type Customer struct {
	UID            string `json:"UID"`
	CreatedDate    time.Time
	Name           string `json:"Name"`
	Address        string `json:"address"`
	Email          string `json:"Contact"`
	PhoneNumber    int    `json:"PhoneNumber"`
	ModifiedDate   time.Time
	RegisteredDate time.Time
}

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
