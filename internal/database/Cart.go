package database

import (
	"log"
)

type CartContent struct {

	//Auto Incrment
	Order int `json:"Order"`
	//From CartRequestPayload
	ItemID   string `json:"ItemID"`
	Quantity int    `json:"Quantity"`
	//From Item
	Status string `json:"Status"`
	//From Item
	ItemName string `json:"ItemName"`
	Price    int    `json:"Price"`
	Stock    int    `json:"Stock"`
}
type CartContents []CartContent

// カートの取得
func GetCartContents(CartID string) (CartContents CartContents, err error) {
	db := ConnectSQL()
	defer db.Close()
	rows, err := db.Query(`
	SELECT 
		Item.Status ,
		Item.Price , 
		Item.Name , 
		Item.Stock ,
		CartContents.CartOrder , 
		CartContents.ItemID , 
		CartContents.Quantity
	FROM 
		Item
	JOIN 
		CartContents 
	ON 
	CartContents.ItemID = Item.ItemID 
	WHERE 
		CartID = ?
	Order By CartContents.CartOrder
		`, CartID)
	if err != nil {
		return CartContents, err
	}
	defer rows.Close()
	for rows.Next() {
		CartContent := new(CartContent)
		err := rows.Scan(&CartContent.Status, &CartContent.Price, &CartContent.ItemName, &CartContent.Stock, &CartContent.Order, &CartContent.ItemID, &CartContent.Quantity)
		if err != nil {
			return CartContents, err
		}
		if CartContent.Quantity > CartContent.Stock {
			CartContent.Status = "OutOfStock"
		}
		CartContents = append(CartContents, *CartContent)
	}
	return CartContents, nil
}

type CartContentRequestPayload struct {
	ItemID   string `json:"ItemId"`
	Quantity int    `json:"Quantity"`
}

// カートの追加・変更・削除
func (c *CartContentRequestPayload) Cart(CartID string) error {
	log.Println("CartID : " + CartID)
	log.Print("ItemID : "+c.ItemID, " Quantity : ", c.Quantity)
	Item := new(Item)
	Item.ItemGet(c.ItemID)
	//リクエストされた商品が登録可能か判定
	log.Println("ItemStatus : " + Item.Status)
	if Item.Status == "Available" {
		Carts, err := GetCartContents(CartID)
		if err != nil {
			return err
		}
		//リクエストされた商品がカートに存在するか確認
		//存在する場合
		if cartContentsSearch(Carts, c.ItemID) {
			if c.Quantity == 0 {
				CartContentDelete(CartID, c.ItemID)
			} else {
				if Item.Stock >= c.Quantity {
					c.CartContentUpdate(CartID)
				} else {
					log.Print("stock is not enough")
				}
			}
		} else {
			//存在しない場合
			if c.Quantity != 0 {
				if Item.Stock >= c.Quantity {
					c.CartContentPost(CartID)
				} else {
					log.Print("stock is not enough")
				}

			} else {
				log.Print("CartReq Quantity is 0")
			}
		}
	}
	return nil
}

// カートに商品を追加
func (c CartContentRequestPayload) CartContentPost(CartID string) error {
	log.Print("CartContentPost")
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの準備
	//UID,ItemID,Quantity
	ins, err := db.Prepare(`
	INSERT 
	INTO 
		CartContents 
		(CartID , ItemID , Quantity) 
		VALUES 
		(? , ? , ?)`)
	log.Print("ins : ", ins, " err : ", err)
	if err != nil {
		return err
	}
	defer ins.Close()
	// SQLの実行
	res, err := ins.Exec(CartID, c.ItemID, c.Quantity)
	log.Print("res : ", res, " err : ", err)
	if err != nil {
		return err
	}
	return nil
}

// カートの商品の変更
func (c CartContentRequestPayload) CartContentUpdate(CartID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare(`
	UPDATE 
		CartContents 
	SET 
		Quantity = ? 
	WHERE 
		CartId = ? 
	AND 
		ItemID = ?`)
	if err != nil {
		return err
	}
	// SQLの実行
	_, err = ins.Exec(c.Quantity, CartID, c.ItemID)
	defer ins.Close()
	if err != nil {
		return err
	}
	return nil
}

// カートから商品を削除
func CartContentDelete(CartID string, ItemID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare(`
	DELETE 
	FROM
		CartContents

	WHERE 
		CartID = ? 

	AND 
		ItemID = ?`)
	if err != nil {
		log.Print(err)
		return err
	}
	// SQLの実行
	_, err = ins.Exec(CartID, ItemID)
	if err != nil {
		return err
	}
	defer ins.Close()
	return nil
}

// 商品自体を削除したときにすべてのカートから特定のアイテムを消す　使用しない方針で行く(商品を削除する前に行う必要がある。)
func CartContentItemDelete(ItemID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare(`
	DELETE 
	FROM 
		CartContents 

	WHERE 
		ItemID = ?`)
	if err != nil {
		log.Print(err)
		return err
	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(ItemID)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// カートから商品を一括で削除　transaction時に使用
func CartDelete(CartID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの準備
	del, err := db.Prepare(`
	DELETE 
	FROM 
		CartContents 
	WHERE 
		CartID = ?
	`)
	if err != nil {
		return err
	}
	defer del.Close()

	// SQLの実行
	del.Exec(CartID)
	return nil
}

// カートに登録する商品が存在するかどうか
func cartContentsSearch(CartContents CartContents, ItemID string) bool {
	for _, CartContent := range CartContents {
		if CartContent.ItemID == ItemID {
			return true
		}
	}
	return false
}
