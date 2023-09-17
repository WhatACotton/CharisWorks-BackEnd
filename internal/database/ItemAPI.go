package database

import (
	"github.com/pkg/errors"
)

type APITopItem struct {
	ItemName string `json:"ItemName"`
	Stock    int    `json:"Stock"`
	Order    int    `json:"Order"`
}
type APITopItems []APITopItem

func ItemGetTop() (APITopItems, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			ItemDetails.Name,
			ItemDetails.Stock,
			Item.ItemOrder 

		FROM 
			ItemDetails

		JOIN 
			Item 

		ON 
			Item.InfoID = ItemDetails.InfoID 

		WHERE 
			ItemDetails.Top = 1 

		AND 
			Item.Status = 'Available'

		ORDER BY 
			Item.ItemOrder`)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetTop1")
	}
	var returnItem []APITopItem
	for rows.Next() {
		TopItem := new(APITopItem)
		err := rows.Scan(&TopItem.ItemName, &TopItem.Stock, &TopItem.Order)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetTop2")
		}
		returnItem = append(returnItem, *TopItem)
	}
	return returnItem, nil
}

type APIItem struct {
	ItemID   string `json:"ItemID"`
	ItemName string `json:"ItemName"`
	Stock    int    `json:"Stock"`
	Price    int    `json:"Price"`
	Status   string `json:"Status"`
	Order    int    `json:"Order"`
}
type APIItems []APIItem

func ItemGetALL() (APIItems, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item.ItemID,
			Item.ItemOrder,
			ItemDetails.Name,
			ItemDetails.Stock,
			ItemDetails.Price,
			Item.Status 

		FROM 
			Item 

		JOIN 
			ItemDetails 

		ON 
			Item.InfoID = ItemDetails.InfoID 

		WHERE 
			Item.Status = 'Available'`)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetALL1")
	}
	var returnItem []APIItem
	for rows.Next() {
		Item := new(APIItem)
		err := rows.Scan(&Item.ItemID, &Item.Order, &Item.ItemName, &Item.Stock, &Item.Price, &Item.Status)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetALL2")
		}
		returnItem = append(returnItem, *Item)
	}
	return returnItem, nil
}

func ItemCategoryGet(Category string) (APIItems, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item.ItemOrder,
			ItemDetails.Name,
			ItemDetails.Stock,
			ItemDetails.Price,
			Item.Status 

		FROM 
			Item

		JOIN 
			ItemDetails

		ON 
			Item.InfoID = ItemDetails.InfoID 

		WHERE 
			Item.Status = 'Available' 

		AND 
			ItemDetails.Category = ?`,
		Category)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetItemCategory1")
	}
	var returnItem []APIItem
	for rows.Next() {
		Item := new(APIItem)
		err := rows.Scan(&Item.Order, &Item.ItemName, &Item.Stock, &Item.Price, &Item.Status)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetItemCategory2")
		}
		returnItem = append(returnItem, *Item)
	}
	return returnItem, nil
}

type ItemDetails struct {
	DetailsID   string `json:"DetailsID"`
	Stock       int    `json:"Stock"`
	Price       int    `json:"Price"`
	ItemName    string `json:"ItemName"`
	Description string `json:"Description"`
	MadeBy      string `json:"MadeBy"`
	Color       string `json:"Color"`
	Series      string `json:"Series"`
}

func ItemColorGet(Color string) (APIItems, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item.ItemOrder,
			ItemDetails.Name,
			ItemDetails.Sstock,
			ItemDetails.Price,
			Item.Status 
		
		FROM 
			Item 
		
		JOIN 
			ItemDetails 
		
		ON 
			Item.InfoID = ItemDetails.InfoID 
		
		WHERE 
			Item.Status = 'Available' 
		
		AND 
			ItemDetails.Color = ?`,
		Color)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetItemColor1")
	}
	var returnItem []APIItem
	for rows.Next() {
		Item := new(APIItem)
		err := rows.Scan(&Item.Order, &Item.ItemName, &Item.Stock, &Item.Price, &Item.Status)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetItemColor2")
		}
		returnItem = append(returnItem, *Item)
	}
	return returnItem, nil
}
