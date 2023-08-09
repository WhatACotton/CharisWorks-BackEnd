package models

// Customer関連
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

type CustomerRequestPayload struct {
	UID         string `json:"UID"`
	Email       string `json:"contact"`
	CreatedDate string
}
type CustomerRegisterPayload struct {
	Name        string `json:"Name"`
	Address     string `json:"address"`
	Email       string `json:"Contact"`
	PhoneNumber int    `json:"PhoneNumber"`
}

// Cart関連
type CartList struct {
	CartId     string `json:"CartId"`
	UID        string `json:"UID"`
	SessionKey string `json:"sessionkey"`
	Valid      bool   `json:"Valid"`
}

type Cart struct {
	CartId         string `json:"CartId"`
	ItemId         string `json:"ItemId"`
	InfoId         string `json:"InfoId"`
	Quantity       int    `json:"Quantity"`
	RegisteredDate string `json:"RegisteredDate"`
	Status         string `json:"Status"`
}

type CartRequestPayload struct {
	ItemId   string `json:"ItemId"`
	Quantity int    `json:"quantity"`
}

// Transaction関連
type TransactionRequestPayload struct {
	UID    string `json:"UID"`
	ItemId string `json:"ItemId"`
	Count  int    `json:"count"`
}

type Transaction struct {
	InfoId   string `json:"ItemId"`
	CartId   string `json:"CartId"`
	Quantity int    `json:"Quantity"`
}
type TransactionList struct {
	CartId          string `json:"CartId"`
	UID             string `json:"UID"`
	TransactionDate string `json:"TransactionDate"`
}

type Bill struct {
	CartId          string        `json:"CartId"`
	UID             string        `json:"UID"`
	TransactionDate string        `json:"TransactionDate"`
	TotalPrice      int           `json:"TotalPrice"`
	TotalCount      int           `json:"TotalCount"`
	Transactions    []Transaction `json:"Items"`
}

// Item関連
type Item struct {
	ItemId string `json:"id"`
	InfoId string `json:"infoid"`
}

type ItemInfo struct {
	InfoId      string `json:"infoid"`
	Price       int    `json:"price"`
	Name        string `json:"Name"`
	Stonesize   int    `json:"Stonesize"`
	Minlength   int    `json:"Minlength"`
	Maxlength   int    `json:"Maxlength"`
	Decsription string `json:"Description"`
	Keyword     string `json:"Keyword"`
}
