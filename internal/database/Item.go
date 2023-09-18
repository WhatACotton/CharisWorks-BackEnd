package database

import "github.com/pkg/errors"

// Item関連
type Item struct {
	ItemID      string `json:"id"`
	DetailsID   string `json:"infoid"`
	Status      string `json:"status"`
	ItemName    string `json:"name"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	MadeBy      MadeBy `json:"madeby"`
	ItemOrder   int    `json:"order"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Series      string `json:"series"`
	Size        int    `json:"size"`
}
type MadeBy string
type Items []Item
type TopItem struct {
	ItemName string `json:"ItemName"`
	Stock    int    `json:"Stock"`
	Order    int    `json:"Order"`
}
type TopItems []TopItem

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

		JOIN 
			ItemDetails

		ON 
			Item.DetailsID = ItemDetails.DetailsID 

		WHERE 
			Item.Status = 'Available' 

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
		
		JOIN 
			ItemDetails 
		
		ON 
			Item.DetailsID = ItemDetails.DetailsID 
		
		WHERE 
			Item.Status = 'Available' 
		
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
func (i *Item) ItemGet(ItemID string) {
	i.ItemID = ItemID
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, _ := db.Query(
		`SELECT
			DetailsID,
			Status,
			ItemName,
			Price,
			Stock,
			MadeBy,
			ItemOrder,
			Description,
			Color,
			Series,
			Size
		FROM
			Item
		WHERE
			ItemID = ?`,
		ItemID)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&i.DetailsID, &i.Status, &i.ItemName, &i.Price, &i.Stock, &i.MadeBy, &i.ItemOrder, &i.Description, &i.Color, &i.Series, &i.Size)
	}
}
func ItemCreate(Item Item) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	db.Exec(`
	INSERT INTO 
		Item 
		(DetailsID,Status,ItemName,Price,Stock,MadeBy,ItemOrder,Description,Color,Series,Size) 
	VALUES 
		(?,?,?,?,?,?,?,?,?,?,?)`,
		Item.DetailsID, Item.Status, Item.ItemName, Item.Price, Item.Stock, Item.MadeBy, Item.ItemOrder, Item.Description, Item.Color, Item.Series, Item.Size)
}
func ItemModify(Item Item) {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	db.Exec(`
	UPDATE 
		Item 
	SET 
		DetailsID = ?,
		Status = ?,
		ItemName = ?,
		Price = ?,
		Stock = ?,
		MadeBy = ?,
		ItemOrder = ?,
		Description = ?,
		Color = ?,
		Series = ?,
		Size = ?
	WHERE 
		ItemID = ?`,
		Item.DetailsID, Item.Status, Item.ItemName, Item.Price, Item.Stock, Item.MadeBy, Item.ItemOrder, Item.Description, Item.Color, Item.Series, Item.Size, Item.ItemID)
}
