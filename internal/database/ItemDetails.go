package database

import "github.com/pkg/errors"

func ItemDetailsCreate() {
	// ItemDetailsの作成
}
func ItemDetailsModyfy() {
	// ItemDetailsの編集
}
func ItemDetailsGet(DetailsID string) (ItemDetails, error) {
	ItemDetails := new(ItemDetails)
	ItemDetails.DetailsID = DetailsID
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			Stock,
			Price,
			ItemName,
			Description,
			MadeBy,
			Color,
			Series
		From 
			ItemDetails

		WHERE 
			DetailsID = ?`,
		DetailsID)
	if err != nil {
		return *ItemDetails, errors.Wrap(err, "error in getting TopItem /GetItemDetails1")
	}
	for rows.Next() {
		err := rows.Scan(&ItemDetails.Stock, &ItemDetails.Price, &ItemDetails.ItemName, &ItemDetails.Description, &ItemDetails.MadeBy, &ItemDetails.Color, &ItemDetails.Series)
		if err != nil {
			return *ItemDetails, errors.Wrap(err, "error in scanning CartID /GetItemDetails2")
		}
	}
	return *ItemDetails, nil
}
func ItemDetailsIDGet(ItemID string) (InfoID string, err error) {
	db := ConnectSQL()
	defer db.Close()
	Status := new(string)
	// SQLの実行
	rows, err := db.Query(
		`SELECT 
			InfoID,
			Status

		FROM 
			Item 
		
		WHERE 
			ItemID = ?`,
		ItemID)
	if err != nil {
		return "error", errors.Wrap(err, "error in getting TopItem /GetInfoID1")
	}
	for rows.Next() {
		err := rows.Scan(&InfoID, &Status)
		if err != nil {
			return "error", errors.Wrap(err, "error in scanning CartID /GetInfoID2")
		}
	}
	if *Status == `Available` {
		return InfoID, nil
	} else {
		return "Couldn't get", nil
	}
}
