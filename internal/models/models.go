package models

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
	CreatedDate string
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
	CreatedDate    string
	ModifiedDate   string
	RegisteredDate string
	LastSessionId  string
}
type LogInLog struct {
	UID       string
	LoginId   string
	LoginDate string
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
	TransactionDate string
}

type Cart struct {
	UID            string
	CartId         string
	ItemId         string
	Quantity       int
	RegisteredDate string
}

type CartRequestPayload struct {
	ItemId   string `json:"itemid"`
	Quantity int    `json:"quantity"`
}

type Session struct {
	SessionId string
	Date      string
}
