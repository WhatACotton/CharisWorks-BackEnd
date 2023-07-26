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
