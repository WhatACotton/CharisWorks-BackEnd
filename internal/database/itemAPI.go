package database

import (
	"github.com/pkg/errors"
)

type API_Top_Item struct {
	Item_Name string `json:"Item_Name"`
	Stock     int    `json:"Stock"`
	Order     int    `json:"Order"`
}
type API_Top_Item_List []API_Top_Item

func Get_Top() (API_Top_Item_List, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item_Info.Name,
			Item_Info.Stock,
			Item_List.Item_Order 

		FROM 
			Item_Info

		JOIN 
			Item_List 

		ON 
			Item_List.Info_ID = Item_Info.Info_ID 

		WHERE 
			Item_Info.Top = 1 

		AND 
			Item_List.Status = 'Available'

		ORDER BY 
			Item_List.Item_Order`)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting Top_Item /Get_Top_1")
	}
	var return_Item []API_Top_Item
	for rows.Next() {
		Top_Item := new(API_Top_Item)
		err := rows.Scan(&Top_Item.Item_Name, &Top_Item.Stock, &Top_Item.Order)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning Cart_ID /Get_Top_2")
		}
		return_Item = append(return_Item, *Top_Item)
	}
	return return_Item, nil
}

type API_Item struct {
	Item_ID   string `json:"Item_ID"`
	Item_Name string `json:"Item_Name"`
	Stock     int    `json:"Stock"`
	Price     int    `json:"Price"`
	Status    string `json:"Status"`
	Order     int    `json:"Order"`
}
type API_Item_List []API_Item

func Get_ALL() (API_Item_List, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item_List.Item_ID,
			Item_List.Item_Order,
			Item_Info.Name,
			Item_Info.Stock,
			Item_Info.Price,
			Item_List.Status 

		FROM 
			Item_List 

		JOIN 
			Item_Info 

		ON 
			Item_List.Info_ID = Item_Info.Info_ID 

		WHERE 
			Item_List.Status = 'Available'`)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting Top_Item /Get_ALL_1")
	}
	var return_Item []API_Item
	for rows.Next() {
		Item := new(API_Item)
		err := rows.Scan(&Item.Item_ID, &Item.Order, &Item.Item_Name, &Item.Stock, &Item.Price, &Item.Status)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning Cart_ID /Get_ALL_2")
		}
		return_Item = append(return_Item, *Item)
	}
	return return_Item, nil
}

func Get_Item_Category(Category string) (API_Item_List, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item_List.Item_Order,
			Item_Info.Name,
			Item_Info.Stock,
			Item_Info.Price,
			Item_List.Status 

		FROM 
			Item_List 

		JOIN 
			Item_Info 

		ON 
			Item_List.Info_ID = Item_Info.Info_ID 

		WHERE 
			Item_List.Status = 'Available' 

		AND 
			Item_Info.Category = ?`,
		Category)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting Top_Item /Get_Item_Category_1")
	}
	var return_Item []API_Item
	for rows.Next() {
		Item := new(API_Item)
		err := rows.Scan(&Item.Order, &Item.Item_Name, &Item.Stock, &Item.Price, &Item.Status)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning Cart_ID /Get_Item_Category_2")
		}
		return_Item = append(return_Item, *Item)
	}
	return return_Item, nil
}

func Get_Info_Id(Item_ID string) (Info_ID string, err error) {
	db := ConnectSQL()
	defer db.Close()
	Status := new(string)
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Info_ID,
			Status

		FROM 
			Item_List 
		
		WHERE 
			Item_ID = ?`,
		Item_ID)
	if err != nil {
		return "error", errors.Wrap(err, "error in getting Top_Item /Get_Info_ID_1")
	}
	for rows.Next() {
		err := rows.Scan(&Info_ID, &Status)
		if err != nil {
			return "error", errors.Wrap(err, "error in scanning Cart_ID /Get_Info_ID_2")
		}
	}
	if *Status == `Available` {
		return Info_ID, nil
	} else {
		return "Couldn't get", nil
	}
}

type API_Item_Details struct {
	Info_ID     string `json:"Info_ID"`
	Item_Name   string `json:"Item_Name"`
	Stock       int    `json:"Stock"`
	Price       int    `json:"Price"`
	Color       string `json:"Color"`
	Description string `json:"Description"`
	Category    string `json:"category"`
	Key_Words   string `json:"Key_Words"`
}

func Get_Item_Details(Info_ID string) (API_Item_Details, error) {
	Item_Details := new(API_Item_Details)
	Item_Details.Info_ID = Info_ID
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Name,
			Price,
			Stock,
			Color,
			Category,
			Key_Words,
			Description 

		From 
			Item_Info

		WHERE 
			Info_ID = ?`,
		Info_ID)
	if err != nil {
		return *Item_Details, errors.Wrap(err, "error in getting Top_Item /Get_Item_Details_1")
	}
	for rows.Next() {
		err := rows.Scan(&Item_Details.Item_Name, &Item_Details.Price, &Item_Details.Stock, &Item_Details.Color, &Item_Details.Category, &Item_Details.Key_Words, &Item_Details.Description)
		if err != nil {
			return *Item_Details, errors.Wrap(err, "error in scanning Cart_ID /Get_Item_Details_2")
		}
	}
	return *Item_Details, nil
}
func Get_Item_Color(Color string) (API_Item_List, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item_List.Item_Order,
			Item_Info.Name,
			Item_Info.Stock,
			Item_Info.Price,
			Item_List.Status 
		
		FROM 
			Item_List 
		
		JOIN 
			Item_Info 
		
		ON 
			Item_List.Info_ID = Item_Info.Info_ID 
		
		WHERE 
			Item_List.Status = 'Available' 
		
		AND 
			Item_Info.Color = ?`,
		Color)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting Top_Item /Get_Item_Color_1")
	}
	var return_Item []API_Item
	for rows.Next() {
		Item := new(API_Item)
		err := rows.Scan(&Item.Order, &Item.Item_Name, &Item.Stock, &Item.Price, &Item.Status)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning Cart_ID /Get_Item_Color_2")
		}
		return_Item = append(return_Item, *Item)
	}
	return return_Item, nil
}
