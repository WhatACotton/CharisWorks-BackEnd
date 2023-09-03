package database

import (
	"log"

	"github.com/pkg/errors"
)

type Cart struct {

	//Auto Incrment
	Order int `json:"Order"`
	//From Cart_Request_Payload
	Item_ID  string `json:"Item_ID"`
	Quantity int    `json:"Quantity"`
	//From Item_List
	Info_ID string `json:"Info_ID"`
	Status  string `json:"Status"`
	//From Item
	Item_Name string `json:"Item_Name"`
	Price     int    `json:"Price"`
	Stock     int    `json:"Stock"`
}

func Get_Cart_Info(Cart_ID string) (Carts []Cart, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	Cart := new(Cart)
	// SQLの実行
	rows, err := db.Query("SELECT Item_List.Info_ID , Item_List.Status , Cart.Order , Cart.Item_ID , Cart.Quantity , Item_Info.Price , Item_Info.Name , Item_Info.Stock FROM Item_Info JOIN Item_List ON Item_List.Info_ID = Item_Info.Info_ID JOIN Cart ON Cart.Item_ID = Item_List.Item_ID WHERE Cart_ID = ?", Cart_ID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting prepare Cart_ID /Get_Cart_Info_1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Cart.Info_ID, &Cart.Status, &Cart.Order, &Cart.Item_ID, &Cart.Quantity, &Cart.Price, &Cart.Item_Name, &Cart.Stock)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning Cart_ID /Get_Cart_Info_2")
		}
		Carts = append(Carts, *Cart)
	}
	return Carts, nil
}

type Cart_Request_Payload struct {
	Item_ID  string `json:"ItemId"`
	Quantity int    `json:"quantity"`
}

func Get_Carts(Cart_ID string) (Carts []Cart, err error) {
	Carts, err = Get_Cart_Info(Cart_ID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting Cart_ID /Get_Carts_1")
	}
	for _, _Cart := range Carts {
		err := _Cart.Get_Cart()
		if err != nil {
			return nil, errors.Wrap(err, "error in getting Cart_ID /Get_Carts_2")
		}
		Carts = append(Carts, _Cart)
	}
	return Carts, nil
}

func (c *Cart) Get_Cart() error {
	err := c.Get_Item_Info()
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Get_Cart_1")
	}
	Item_List := new(Item_List)
	err = Item_List.Get_Item_List(c.Item_ID)
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Get_Cart_2")
	}
	c.Status = Item_List.Status
	return nil
}
func (c *Cart) Get_Item_Info() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query("SELECT Item_List.Item_ID,Item_List.Info_ID,Item_List.Status,Item_Info.Info_ID,Item_Info.Name ,Item_Info.Stock FROM Item_List JOIN Item_Info ON Item_List.Info_ID = Item_Info.Info_ID WHERE Item_ID = ?", c.Item_ID)
	if err != nil {
		return errors.Wrap(err, "error in getting prepare Cart_ID /Get_Item_Info_1")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&c.Item_ID, &c.Info_ID, &c.Status)
		if err != nil {
			return errors.Wrap(err, "error in scanning Cart_ID /Get_Item_Info_2")
		}
	}
	return nil
}

func (c Cart_Request_Payload) Cart(Cart_ID string) error {
	log.Println("Cart_ID : " + Cart_ID)
	Item_List := new(Item_List)
	err := Item_List.Get_Item_List(c.Item_ID)
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Cart_1")
	}
	//リクエストされた商品が登録可能か判定
	log.Println("Item_List.Status : " + Item_List.Status)
	if Item_List.Status == "Available" {
		Carts, err := Get_Cart_Info(Cart_ID)
		if err != nil {
			return errors.Wrap(err, "error in getting Cart_ID /Cart_2")
		}
		//リクエストされた商品がカートに存在するか確認
		if Search_Cart(Carts, c.Item_ID) {
			if c.Quantity == 0 {
				err := Delete_Cart(Cart_ID, c.Item_ID)
				if err != nil {
					return errors.Wrap(err, "error in getting Cart_ID /Cart_3")
				}
			} else {

				if Item_List.Stock >= c.Quantity {
					err := c.Update_Cart(Cart_ID)
					if err != nil {
						return errors.Wrap(err, "error in getting Cart_ID /Cart_4")
					}
				} else {
					return errors.New("stock is not enough")
				}
			}
		} else {
			if c.Quantity != 0 {
				if Item_List.Stock >= c.Quantity {
					err := c.Post_Cart(Cart_ID)
					if err != nil {
						return errors.Wrap(err, "error in getting Cart_ID /Cart_5")
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

func (c Cart_Request_Payload) Post_Cart(Cart_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの準備
	//UID,Item_ID,Quantity
	ins, err := db.Prepare("INSERT INTO Cart (Cart_ID , Item_ID , Quantity) VALUES (? , ? , ?)")
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Post_Cart_1")
	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(Cart_ID, c.Item_ID, c.Quantity)
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Post_Cart_2")
	}
	return nil
}

func Search_Cart(Carts []Cart, Item_ID string) bool {
	for _, Cart := range Carts {
		if Cart.Item_ID == Item_ID {
			return true
		}
	}
	return false
}

func Start_Cart(SessionKey string, Cart_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの準備
	//UID,Item_ID,Quantity
	ins, err := db.Prepare("INSERT INTO Cart_List (Cart_ID,Session_Key)VALUES(?,?)")
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Start_Cart_1")
	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(Cart_ID, SessionKey)
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Start_Cart_2")
	}
	return nil
}

func Get_Cart_ID(OldSessionKey string) (Cart_ID string, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT CartId FROM cartlist WHERE SessionKey = ?", OldSessionKey)
	if err != nil {
		return "none", errors.Wrap(err, "error in getting Cart_ID /Get_Cart_ID_1")
	}

	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Cart_ID)

		if err != nil {
			return "none", errors.Wrap(err, "error in scanning Cart_ID /Get_Cart_ID_2")
		}
	}
	return Cart_ID, nil
}

func (c Cart_Request_Payload) Update_Cart(Cart_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare("UPDATE Cart SET Quantity = ? WHERE Cart_Id = ? AND Item_ID = ?")
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Update_Cart_1")
	}
	// SQLの実行
	_, err = ins.Exec(c.Quantity, Cart_ID, c.Item_ID)
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Update_Cart_2")
	}
	defer ins.Close()
	return nil
}

func Delete_Cart(Cart_ID string, Item_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare("DELETE FROM Cart WHERE Cart_ID = ? AND Item_ID = ?")
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Delete_Cart_1")
	}
	// SQLの実行
	_, err = ins.Exec(Cart_ID, Item_ID)
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Delete_Cart_2")
	}
	defer ins.Close()
	return nil
}

func Delete_Item_From_Cart(Item_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare("DELETE FROM Cart WHERE Item_ID = ?")
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Delete_Item_From_Cart_1")
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(Item_ID)
	if err != nil {
		return errors.Wrap(err, "error in getting Cart_ID /Delete_Item_From_Cart_2")
	}
	return nil
}
