package models

type LogInLog struct {
	UserID    string
	LoginId   string
	LoginDate string
}

// Transaction関連
type TransactionRequestPayload struct {
	UserID string `json:"UserID"`
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
	UserID          string `json:"UserID"`
	TransactionDate string `json:"TransactionDate"`
}

type Bill struct {
	CartId          string        `json:"CartId"`
	UserID          string        `json:"UserID"`
	Name            string        `json:"Name"`
	Address         string        `json:"address"`
	PhoneNumber     string        `json:"PhoneNumber"`
	Email           string        `json:"Contact"`
	TransactionDate string        `json:"TransactionDate"`
	TotalPrice      int           `json:"TotalPrice"`
	TotalCount      int           `json:"TotalCount"`
	Transactions    []Transaction `json:"Items"`
}
