package models

type CustomerRequestPayload struct {
	UID         string `json:"uid"`
	CreatedDate string `json:"CreatedDate"`
	Contact     string `json:"contact"`
}
type Customer struct {
	UID         string `json:"UID"`
	CreatedDate string `json:"CreatedDate"`
	Name        string `json:"Name"`
	Address     string `json:"address"`
	Contact     string `json:"Contact"`
}

type TransactionRequestPayload struct {
	UID        string `json:"UID"`
	ItemId     string `json:"itemid"`
	Count      int    `json:"count"`
	IsFinished bool   `json:"isFinished"`
}
type Transaction struct {
	UID             string `json:"UID"`
	ItemId          string `json:"itemid"`
	TransactionId   string `json:"transactionId"`
	TransactionDate string `json:"transactionDate"`
	Count           string `json:"count"`
	IsFinished      bool   `json:"isFinished"`
}

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
type PatchRequestPayload struct {
	ID        string `json:"id"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
	Isint     bool   `json:"isint"`
}
