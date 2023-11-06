package database

import (
	"log"

	"github.com/WhatACotton/go-backend-test/validation"
	"github.com/pkg/errors"
)

// Item関連
type Item struct {
	ItemID          string `json:"ItemID"`
	StripeAccountID string `json:"StripeAccountID"`
	Status          string `json:"Status"`
	Name            string `json:"Name"`
	Price           int    `json:"Price"`
	Stock           int    `json:"Stock"`
	Order           int    `json:"Order"`
	Color           string `json:"Color"`
	Series          string `json:"Series"`
	Size            string `json:"Size"`
	Description     string `json:"Description"`
	//MakersDetailsから取得
	MakerName        string `json:"MakerName"`
	MakerDescription string `json:"MakerDescription,omitempty"`
}
type ItemMain struct {
	ItemID      string `json:"ItemID"`
	Status      string `json:"Status"`
	Price       int    `json:"Price"`
	Stock       int    `json:"Stock"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Color       string `json:"Color"`
	Series      string `json:"Series"`
	Size        string `json:"Size"`
}
type ItemDetail struct {
	ItemID      string `json:"ItemID"`
	Description string `json:"Description"`
	Color       string `json:"Color"`
	Series      string `json:"Series"`
	Size        string `json:"Size"`
}
type ItemForList struct {
	ItemID string            `json:"ItemID"`
	Name   string            `json:"Name"`
	Price  int               `json:"Price"`
	Stock  int               `json:"Stock"`
	Order  int               `json:"Order"`
	Tags   map[string]string `json:"tags"`
}
type Items []Item
type TopItem struct {
	Name  string `json:"Name"`
	Stock int    `json:"Stock"`
}
type TopItems []TopItem

func IsItemExist(ItemID string) (bool, int) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, _ := db.Query(
		`SELECT
			ItemID,
			Stock
		FROM 
			Item 
		WHERE 
			ItemID = ?`,
		ItemID)
	defer rows.Close()
	Stock := 0
	for rows.Next() {
		rows.Scan(&ItemID, &Stock)
	}
	if ItemID == "" {
		return false, 0
	} else {
		return true, Stock
	}
}

// Itemの取得
func (i *Item) ItemGet(ItemID string) {
	i.ItemID = ItemID
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, _ := db.Query(
		`SELECT
			Item.Status,
			Item.ItemName,
			Item.Price,
			Item.Stock,
			Item.ItemOrder,
			Item.Color,
			Item.Series,
			Item.Size,
			Item.Description,
			Maker.MakerName,
			Maker.Description
		FROM 
			Item 
		JOIN 
			Maker
		ON 
			Item.StripeAccountID = Maker.StripeAccountID 
		WHERE 
			ItemID = ?`,
		ItemID)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&i.Status, &i.Name, &i.Price, &i.Stock, &i.Order, &i.Color, &i.Series, &i.Size, &i.Description, &i.MakerName, &i.MakerDescription)
	}
}

// Top用Itemの取得
func ItemGetTop() (TopItems, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item.ItemName,
			Item.Stock,
			Item.ItemOrder 
		FROM 
			Item
		WHERE
			Item.Status = 'Available'
			AND
			TOP = 1
		ORDER BY 
			Item.Order`)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetTop1")
	}
	var returnItem []TopItem
	for rows.Next() {
		TopItem := new(TopItem)
		err := rows.Scan(&TopItem.Name, &TopItem.Stock)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetTop2")
		}
		returnItem = append(returnItem, *TopItem)
	}
	return returnItem, nil
}

// すべてのItemの取得
func ItemGetALL() (Items Items, err error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item.ItemID,
			Item.ItemOrder,
			Item.Status,
			Item.Price,
			Item.Stock,
			Item.ItemName,
			Item.Description,
			Item.Color,
			Item.Series,
			Item.Size
		FROM 
			Item 
		WHERE 
			Item.Status = 'Available'`)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetALL1")
	}
	for rows.Next() {
		Item := new(Item)
		err := rows.Scan(&Item.ItemID, &Item.Order, &Item.Status, &Item.Price, &Item.Stock, &Item.Name, &Item.Description, &Item.Color, &Item.Series, &Item.Size)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetALL2")
		}
		Items = append(Items, *Item)
	}
	return Items, nil
}

// カテゴリごとのItemの取得
func ItemGetCategory(Category string) (Items, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item.Order,
			Item.ItemName,
			Item.Stock,
			Item.Price,
			Item.Status 

		FROM 
			Item
		AND 
			ItemDetails.Category = ?`,
		Category)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetItemCategory1")
	}
	var returnItem []Item
	for rows.Next() {
		Item := new(Item)
		err := rows.Scan(&Item.Order, &Item.Name, &Item.Stock, &Item.Price, &Item.Status)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetItemCategory2")
		}
		returnItem = append(returnItem, *Item)
	}
	return returnItem, nil
}

// カラーごとのItemの取得
func ItemGetColor(Color string) (Items, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			ItemOrder,
			ItemName,
			Stock,
			Price,
			Status 
		FROM 
			Item
		WHERE 
			Color = ?`,
		Color)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetItemColor1")
	}
	var returnItem []Item
	for rows.Next() {
		Item := new(Item)
		err := rows.Scan(&Item.Order, &Item.Name, &Item.Stock, &Item.Price, &Item.Status)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetItemColor2")
		}
		returnItem = append(returnItem, *Item)
	}
	return returnItem, nil
}

// 出品者ごとのItemの取得
func ItemGetMaker(StribeAccountID string) (Items Items) {
	db := ConnectSQL()
	rows, _ := db.Query(
		`SELECT 
			ItemID,
			Status,
			ItemName,
			Price,
			Stock,
			ItemOrder,
			Color,
			Series,
			Size,
			Description
		FROM 
			Item 
		WHERE 
			StripeAccountID = ?`,
		StribeAccountID)
	defer db.Close()
	for rows.Next() {
		Item := new(Item)
		rows.Scan(&Item.ItemID, &Item.Status, &Item.Name, &Item.Price, &Item.Stock, &Item.Order, &Item.Color, &Item.Series, &Item.Size, &Item.Description)
		Items = append(Items, *Item)
	}
	return Items
}

// Itemの主要情報の作成
func ItemMainCreate(ItemMain ItemMain, StripeAccountID string) {
	ItemID := validation.GetUUID()
	log.Print("MakerName:", StripeAccountID)
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	db.Exec(`
	INSERT INTO 
		Item 
		(ItemID,Status,ItemName,Price,Stock,StripeAccountID,Description,Color,Series,Size) 
	VALUES 
		(?,?,?,?,?,?,?,?,?,?)`,
		ItemID, ItemMain.Status, ItemMain.Name, ItemMain.Price, ItemMain.Stock, StripeAccountID, ItemMain.Description, ItemMain.Color, ItemMain.Series, ItemMain.Size)
}

// Itemの詳細の作成
func ItemDetailCreate(ItemDetail ItemDetail, StripeAccountID string) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	db.Exec(`
	UPDATE
		Item
	SET
		Description = ?,
		Color = ?,
		Series = ?,
		Size = ?
	WHERE
		ItemID = ?
		AND
		StripeAccountID = ?`,
		ItemDetail.Description, ItemDetail.Color, ItemDetail.Series, ItemDetail.Size, ItemDetail.ItemID, StripeAccountID)
}
func CartDetails(ItemID string) (CartContent CartContent) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, _ := db.Query(
		`SELECT
			ItemID,
			ItemName,
			Price,
			Status,
			Stock
		FROM 
			Item  
		WHERE 
			ItemID = ?`,
		ItemID)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&CartContent.ItemID, &CartContent.Name, &CartContent.Price, &CartContent.Status, &CartContent.Stock)
	}
	return CartContent
}
