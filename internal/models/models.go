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
	PhoneNumber    string `json:"PhoneNumber"`
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
	Available bool
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

type CartList struct {
	CartId         string
	UID            string
	SessionKey     string
	RegisteredDate string
	Valid          bool
}

type Cart struct {
	CartId   string
	ItemId   string
	Quantity int
}

type CartRequestPayload struct {
	ItemId   string `json:"itemid"`
	Quantity int    `json:"quantity"`
}

type CartSession struct {
	CartId    string
	SessionId string
	Date      string
	Available bool
}
