package database

import (
	"log"

	"github.com/pkg/errors"
)

// Item関連
type Item struct {
	ItemID      string `json:"ItemID"`
	Status      string `json:"Status"`
	Name        string `json:"Name"`
	Price       int    `json:"price"`
	Stock       int    `json:"Stock"`
	MakerName   string `json:"MakerName"`
	Order       int    `json:"Order"`
	Color       string `json:"Color"`
	Series      string `json:"Series"`
	Size        string `json:"Size"`
	Description string `json:"Description"`
	//MakersDetailsから取得
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

// Itemの取得
func (i *Item) ItemGet(ItemID string) {
	i.ItemID = ItemID
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT
			Item.Status,
			Item.Name,
			Item.Price,
			Item.Stock,
			Item.MakerName,
			Item.ItemOrder,
			Item.Color,
			Item.Series,
			Item.Size,
			Customer.MakerDescription
		FROM 
			Item 
		JOIN 
			Customer 
		ON 
			Item.MakerName = Customer.MakerName 
		WHERE 
			ItemID = ?`,
		ItemID)
	log.Print("rows:", rows, "err:", err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&i.Status, &i.Name, &i.Price, &i.Stock, &i.MakerName, &i.Order, &i.Color, &i.Series, &i.Size, &i.MakerDescription)
		log.Print("Item:", i, "err:", err)
	}
}

// Top用Itemの取得
func ItemGetTop() (TopItems, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item.Name,
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
			Item.Name,
			Item.Description,
			Item.Color,
			Item.Series,
			Item.Size,
			Item.MakerName
		FROM 
			Item 
		WHERE 
			Item.Status = 'Available'`)
	log.Print("rows:", rows, "err:", err)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetALL1")
	}
	for rows.Next() {
		Item := new(Item)
		err := rows.Scan(&Item.ItemID, &Item.Order, &Item.Status, &Item.Price, &Item.Stock, &Item.Name, &Item.Description, &Item.Color, &Item.Series, &Item.Size, &Item.MakerName)
		log.Print("Item:", Item, "err:", err)
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
			Item.Name,
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
			Name,
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
func ItemGetMaker(MakerName string) (Items Items) {
	db := ConnectSQL()
	rows, err := db.Query(
		`SELECT 
			ItemID,
			Status,
			Name,
			Price,
			Stock,
			MakerName,
			ItemOrder,
			Color,
			Series,
			Size,
			Description
		FROM 
			Item 
		WHERE 
			MakerName = ?`,
		MakerName)
	defer db.Close()
	log.Print("rows:", rows, "err:", err)
	for rows.Next() {
		Item := new(Item)
		err := rows.Scan(&Item.ItemID, &Item.Status, &Item.Name, &Item.Price, &Item.Stock, &Item.MakerName, &Item.Order, &Item.Color, &Item.Series, &Item.Size, &Item.Description)
		log.Print("Item:", Item, "err:", err)
		Items = append(Items, *Item)
	}
	return Items
}

// Itemの主要情報の作成
func ItemMainCreate(ItemMain ItemMain, MakerName string) {
	log.Print("MakerName:", MakerName)
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	res, err := db.Exec(`
	INSERT INTO 
		Item 
		(ItemID,Status,Name,Price,Stock,MakerName,Description,Color,Series,Size) 
	VALUES 
		(?,?,?,?,?,?)`,
		ItemMain.ItemID, ItemMain.Status, ItemMain.Name, ItemMain.Price, ItemMain.Stock, MakerName, ItemMain.Description, ItemMain.Color, ItemMain.Series, ItemMain.Size)
	log.Print("res:", res, "err:", err)
}

// Itemの詳細の作成
func ItemDetailCreate(ItemDetail ItemDetail, MakerName string) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	res, err := db.Exec(`
	UPDATE
		Item
	SET
		Description = ?,
		Color = ?,
		Series = ?,
		Size = ?
	WHERE
		ItemID = ?`,
		ItemDetail.Description, ItemDetail.Color, ItemDetail.Series, ItemDetail.Size, ItemDetail.ItemID)
	log.Print("res:", res, "err:", err)
}
