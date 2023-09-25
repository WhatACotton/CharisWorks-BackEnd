package database

import (
	"github.com/pkg/errors"
)

// Item関連
type Item struct {
	ItemID    string `json:"id"`
	Status    string `json:"status"`
	ItemName  string `json:"name"`
	Price     int    `json:"price"`
	Stock     int    `json:"stock"`
	MakerName string `json:"makername"`
	ItemOrder int    `json:"order"`
	Color     string `json:"color"`
	Series    string `json:"series"`
	Size      int    `json:"size"`
	//MakersDetailsから取得
	Description string `json:"description"`
}
type ItemMain struct {
	ItemID   string `json:"ItemID"`
	Status   string `json:"Status"`
	Price    int    `json:"Price"`
	Stock    int    `json:"Stock"`
	ItemName string `json:"ItemName"`
}
type ItemDetail struct {
	ItemID      string `json:"ItemID"`
	Description string `json:"Description"`
	Color       string `json:"Color"`
	Series      string `json:"Series"`
	Size        string `json:"Size"`
}
type ItemForList struct {
	ItemID    string            `json:"id"`
	ItemName  string            `json:"name"`
	Price     int               `json:"price"`
	Stock     int               `json:"stock"`
	ItemOrder int               `json:"order"`
	Tags      map[string]string `json:"tags"`
}
type Items []Item
type TopItem struct {
	ItemName string `json:"ItemName"`
	Stock    int    `json:"Stock"`
	Order    int    `json:"Order"`
}
type TopItems []TopItem

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
			Item.MakerName,
			Item.ItemOrder,
			Item.Color,
			Item.Series,
			Item.Size,
			MakersDetails.Description,
		FROM
			Item
		JOIN
			MakersDetails
		ON
			Item.MakerName = MakersDetails.MakerName
		WHERE
			ItemID = ?`,
		ItemID)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&i.Status, &i.ItemName, &i.Price, &i.Stock, &i.MakerName, &i.ItemOrder, &i.Color, &i.Series, &i.Size, &i.Description)
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
			TOP = 'true'
		ORDER BY 
			ItemDetails.ItemOrder`)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetTop1")
	}
	var returnItem []TopItem
	for rows.Next() {
		TopItem := new(TopItem)
		err := rows.Scan(&TopItem.ItemName, &TopItem.Stock, &TopItem.Order)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetTop2")
		}
		returnItem = append(returnItem, *TopItem)
	}
	return returnItem, nil
}

// すべてのItemの取得
func ItemGetALL() (Items, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item.ItemOrder,
			Item.ItemName,
			Item.Stock,
			Item.Price,
			Item.Status 

		FROM 
			Item 

		WHERE 
			Item.Status = 'Available'`)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetALL1")
	}
	var returnItem Items
	for rows.Next() {
		Item := new(Item)
		err := rows.Scan(&Item.ItemOrder, &Item.ItemName, &Item.Stock, &Item.Price, &Item.Status)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetALL2")
		}
		returnItem = append(returnItem, *Item)
	}
	return returnItem, nil
}

// カテゴリごとのItemの取得
func ItemGetCategory(Category string) (Items, error) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Item.ItemOrder,
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
		err := rows.Scan(&Item.ItemOrder, &Item.ItemName, &Item.Stock, &Item.Price, &Item.Status)
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
			Item.ItemOrder,
			Item.ItemName,
			Item.Stock,
			Item.Price,
			Item.Status 
		
		FROM 
			Item
		AND 
			ItemDetails.Color = ?`,
		Color)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting TopItem /GetItemColor1")
	}
	var returnItem []Item
	for rows.Next() {
		Item := new(Item)
		err := rows.Scan(&Item.ItemOrder, &Item.ItemName, &Item.Stock, &Item.Price, &Item.Status)
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
	rows, _ := db.Query(
		`SELECT 
			Item.ItemID,
			Item.Status,
			Item.ItemName,
			Item.Price,
			Item.Stock,
			Item.ItemOrder

		FROM 
			Item
		WHERE
			MakerName = ?`,
		MakerName)
	defer db.Close()
	for rows.Next() {
		Item := new(Item)
		rows.Scan(&Item.ItemID, &Item.Status, &Item.ItemName, &Item.Price, &Item.Stock, &Item.ItemOrder)
		Items = append(Items, *Item)
	}
	return Items
}

// Itemの主要情報の作成
func ItemMainCreate(ItemMain ItemMain, MakerName string) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	db.Exec(`
	INSERT INTO 
		Item 
		(ItemID,Status,ItemName,Price,Stock,MakerName) 
	VALUES 
		(?,?,?,?,?,?)`,
		ItemMain.ItemID, ItemMain.Status, ItemMain.ItemName, ItemMain.Price, ItemMain.Stock, MakerName)
}

// Itemの詳細の作成
func ItemDetailCreate(ItemDetail ItemDetail, MakerName string) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	db.Exec(`
	INSERT INTO 
		ItemDetails 
		(ItemID,Description,Color,Series,Size) 
	VALUES 
		(?,?,?,?,?)`,
		ItemDetail.ItemID, ItemDetail.Description, ItemDetail.Color, ItemDetail.Series, ItemDetail.Size)
}
