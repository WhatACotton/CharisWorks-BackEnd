package database

import (
	"log"

	"github.com/pkg/errors"
)

type CartContent struct {

	//Auto Incrment
	Order int `json:"Order"`
	//From CartRequestPayload
	ItemID   string `json:"ItemID"`
	Quantity int    `json:"Quantity"`
	//From ItemList
	InfoID string `json:"InfoID"`
	Status string `json:"Status"`
	//From Item
	ItemName string `json:"ItemName"`
	Price    int    `json:"Price"`
	Stock    int    `json:"Stock"`
}
type CartContents []CartContent

func GetCartContents(CartID string) (CartContents CartContents, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		ItemDetails.InfoID , 
		Item.Status ,
		CartContent.Order , 
		CartContent.ItemID , 
		CartContent.Quantity , 
		ItemDetails.Price , 
		ItemDetails.Name , 
		ItemDetails.Stock 
	
	FROM 
		ItemDetails
	
	JOIN 
		Item

	ON 
		Item.InfoID = ItemDetails.InfoID 
	
	JOIN 
		CartContent 

	ON 
	CartContent.ItemID = Item.ItemID 

	WHERE 
		CartID = ?`, CartID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting prepare CartID /GetCartInfo1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		CartContent := new(CartContent)
		err := rows.Scan(&CartContent.InfoID, &CartContent.Status, &CartContent.Order, &CartContent.ItemID, &CartContent.Quantity, &CartContent.Price, &CartContent.ItemName, &CartContent.Stock)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning CartID /GetCartInfo2")
		}
		CartContents = append(CartContents, *CartContent)
	}
	return CartContents, nil
}

type CartRequestPayload struct {
	ItemID   string `json:"ItemId"`
	Quantity int    `json:"Quantity"`
}

func (c *CartRequestPayload) Cart(CartID string) error {
	log.Println("CartID : " + CartID)
	log.Print("ItemID : "+c.ItemID, " Quantity : ", c.Quantity)
	ItemList := new(ItemList)
	err := ItemList.GetItemList(c.ItemID)
	if err != nil {
		return errors.Wrap(err, "error in getting CartID /Cart1")
	}
	//リクエストされた商品が登録可能か判定
	log.Println("ItemStatus : " + ItemList.Status)
	if ItemList.Status == "Available" {
		Carts, err := GetCartContents(CartID)
		if err != nil {
			return errors.Wrap(err, "error in getting CartID /Cart2")
		}
		//リクエストされた商品がカートに存在するか確認
		if SearchCart(Carts, c.ItemID) {
			if c.Quantity == 0 {
				err := DeleteCartContent(CartID, c.ItemID)
				if err != nil {
					return errors.Wrap(err, "error in getting CartID /Cart3")
				}
			} else {

				if ItemList.Stock >= c.Quantity {
					err := c.UpdateCart(CartID)
					if err != nil {
						return errors.Wrap(err, "error in getting CartID /Cart4")
					}
				} else {
					return errors.New("stock is not enough")
				}
			}
		} else {
			if c.Quantity != 0 {
				if ItemList.Stock >= c.Quantity {
					err := c.PostCart(CartID)
					if err != nil {
						return errors.Wrap(err, "error in getting CartID /Cart5")
					}
				} else {
					return errors.New("stock is not enough")
				}

			} else {
				log.Println("CartReq Quantity is 0")
			}
		}
	}
	return nil
}

func (c CartRequestPayload) PostCart(CartID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの準備
	//UID,ItemID,Quantity
	ins, err := db.Prepare(`
	INSERT 
	INTO 
		CartContent 
		(CartID , ItemID , Quantity) 
		VALUES 
		(? , ? , ?)`)
	if err != nil {
		return errors.Wrap(err, "error in getting CartID /PostCart1")
	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(CartID, c.ItemID, c.Quantity)
	if err != nil {
		return errors.Wrap(err, "error in getting CartID /PostCart2")
	}
	return nil
}

func SearchCart(CartContents CartContents, ItemID string) bool {
	for _, CartContent := range CartContents {
		if CartContent.ItemID == ItemID {
			return true
		}
	}
	return false
}

func GetCartID(UID string) (string, error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query(`
	SELECT 
		CartID
	FROM 
		Customer 
	WHERE 
		UID = ?`, UID)
	if err != nil {
		return "none", errors.Wrap(err, "error in getting CartID /GetCartID1")
	}

	defer rows.Close()
	CartID := new(string)
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&CartID)

		if err != nil {
			return "none", errors.Wrap(err, "error in scanning CartID /GetCartID2")
		}
	}
	return *CartID, nil
}

func (c CartRequestPayload) UpdateCart(CartID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare(`
	UPDATE 
		CartContent 
	SET 
		Quantity = ? 
	WHERE 
		CartId = ? 
	AND 
		ItemID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in getting CartID /UpdateCart1")
	}
	// SQLの実行
	_, err = ins.Exec(c.Quantity, CartID, c.ItemID)
	if err != nil {
		return errors.Wrap(err, "error in getting CartID /UpdateCart2")
	}
	defer ins.Close()
	return nil
}

func DeleteCartContent(CartID string, ItemID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare(`
	DELETE 
	FROM
		CartContent

	WHERE 
		CartID = ? 

	AND 
		ItemID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in getting CartID /DeleteCart1")
	}
	// SQLの実行
	_, err = ins.Exec(CartID, ItemID)
	if err != nil {
		return errors.Wrap(err, "error in getting CartID /DeleteCart2")
	}
	defer ins.Close()
	return nil
}

func DeleteCartforTransaction(CartID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare(`
	DELETE 
	FROM
		CartContent

	WHERE 
		CartID = ? 
		`)
	if err != nil {
		return errors.Wrap(err, "error in getting CartID /DeleteCart1")
	}
	// SQLの実行
	_, err = ins.Exec(CartID)
	if err != nil {
		return errors.Wrap(err, "error in getting CartID /DeleteCart2")
	}
	defer ins.Close()
	return nil
}

func DeleteItemFromCart(ItemID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare(`
	DELETE 
	FROM 
		CartContent 

	WHERE 
		ItemID = ?`)
	if err != nil {
		return errors.Wrap(err, "error in getting CartID /DeleteItemFromCart1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(ItemID)
	if err != nil {
		return errors.Wrap(err, "error in getting CartID /DeleteItemFromCart2")
	}
	return nil
}
