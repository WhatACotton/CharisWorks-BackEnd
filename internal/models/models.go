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
	Email       string `json:"contact"`
	CreatedDate time.Time
}
type CustomerRegisterPayload struct {
	Name        string `json:"Name"`
	Address     string `json:"address"`
	Email       string `json:"Contact"`
	PhoneNumber int    `json:"PhoneNumber"`
}
type Customer struct {
	UID            string `json:"UID"`
	Name           string `json:"Name"`
	Address        string `json:"address"`
	Email          string `json:"Contact"`
	PhoneNumber    int    `json:"PhoneNumber"`
	Register       bool
	CreatedDate    time.Time
	ModifiedDate   time.Time
	RegisteredDate time.Time
	LastLogInDate  time.Time
}
type LogInLog struct {
	UID     string
	LoginId string
	Login   time.Time
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
	Count           string `json:"count"`
	IsFinished      bool   `json:"isFinished"`
	TransactionDate time.Time
}

type Cart struct {
	UID            string
	ItemId         string
	Quantity       int
	CartId         string
	RegisteredDate time.Time
}

type CartRequestPayload struct {
	ItemId   string
	Quantity int
}

type Session struct {
	SessionId string
	Date      time.Time
}
