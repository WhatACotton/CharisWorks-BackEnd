package database

import (
	"fmt"
)

type Cart struct {
	Cart_ID   string `json:"CartId"`
	Info_ID   string `json:"InfoId"`
	Item_ID   string `json:"ItemId"`
	Quantity  int    `json:"Quantity"`
	Order     int    `json:"Order"`
	Status    string `json:"Status"`
	Item_Name string `json:"ItemName"`
	Price     int    `json:"Price"`
}

type Return_Cart struct {
	Quantity  int    `json:"Quantity"`
	Order     int    `json:"Order"`
	Status    string `json:"Status"`
	Item_Name string `json:"ItemName"`
	Price     int    `json:"Price"`
}

type Cart_Request_Payload struct {
	Item_ID  string `json:"ItemId"`
	Quantity int    `json:"quantity"`
}

func Get_Return_Cart(Carts []Cart) (Return_Carts []Return_Cart) {
	for _, _Cart := range Carts {
		Return_Cart := new(Return_Cart)
		Return_Cart.Quantity = _Cart.Quantity
		Return_Cart.Order = _Cart.Order
		Return_Cart.Status = _Cart.Status
		Return_Cart.Item_Name = _Cart.Item_Name
		Return_Cart.Price = _Cart.Price
		Return_Carts = append(Return_Carts, *Return_Cart)
	}
	return Return_Carts
}
func Get_Cart(Cart_ID string) (Carts []Cart, err error) {
	Carts, err = Get_Cart_Info(Cart_ID)
	if err != nil {
		return nil, err
	}
	Carts = Get_Cart_With_Info(Carts)
	for _, _Cart := range Carts {
		Cart := new(Cart)
		err := Cart.Get_Item_Name_And_Price()
		if err != nil {
			return nil, err
		}
		Cart.Quantity = _Cart.Quantity
		Cart.Order = _Cart.Order
		Cart.Status = _Cart.Status
		Cart.Item_Name = _Cart.Item_Name
		Cart.Price = _Cart.Price
		Carts = append(Carts, *Cart)
	}
	return Carts, nil
}
func Get_Cart_Info(Cart_ID string) (Carts []Cart, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	Cart := new(Cart)
	// SQLの実行
	rows, err := db.Query("SELECT * FROM Cart WHERE Cart_ID = ?", Cart_ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Cart.Cart_ID, &Cart.Item_ID, &Cart.Quantity)
		if err != nil {
			return nil, err
		}
		Carts = append(Carts, *Cart)
	}
	return Carts, nil
}

func (c Cart_Request_Payload) Cart(Cart_ID string) error {
	Item_List := GetItemList()
	//リクエストされた商品が存在するか確認
	if Inspect_Items(c.Item_ID, Item_List) {
		Carts, err := Get_Cart_Info(Cart_ID)
		if err != nil {
			return err
		}
		//リクエストされた商品がカートに存在するか確認
		if Search_Cart(Carts, c.Item_ID) {
			if c.Quantity == 0 {
				err := Delete_Cart(Cart_ID, c.Item_ID)
				if err != nil {
					return err
				}
			} else {
				err := c.Update_Cart(Cart_ID)
				if err != nil {
					return err
				}
			}
		} else {
			err := c.Post_Cart(Cart_ID)
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		return fmt.Errorf("cart: item not found")
	}
}

func (c Cart_Request_Payload) Post_Cart(Cart_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの準備
	//UID,Item_ID,Quantity
	ins, err := db.Prepare("INSERT INTO cart (Cart_ID , Item_ID , Quantity) VALUES (? , ? , ?)")
	if err != nil {
		return err
	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(Cart_ID, c.Item_ID, c.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func Inspect_Items(Item_ID string, Itemlist []Item_List) bool {
	for _, Item := range Itemlist {
		if Item_ID == Item.Item_ID {
			return true
		}
	}
	return false
}

func Search_Cart(Carts []Cart, Item_ID string) bool {
	for _, Cart := range Carts {
		if Cart.Item_ID == Item_ID {
			return true
		} else {
			return false
		}
	}
	return false
}

func Get_Cart_With_Info(Carts []Cart) []Cart {
	for _, Cart := range Carts {
		Cart.Get_Cart_Info()
	}
	return Carts
}

func (c *Cart) Get_Item_Name_And_Price() error {
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query("SELECT Item_Name , Price FROM Item WHERE Info_ID = ?", c.Info_ID)
	if err != nil {
		return err
	}
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&c.Item_Name, &c.Price)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *Cart) Get_Cart_Info() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query("SELECT * FROM Item_List WHERE Item_ID = ?", c.Item_ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&c.Item_ID, &c.Info_ID, &c.Status)
		if err != nil {
			return err
		}
	}
	return nil
}

func Start_Cart(SessionKey string, Cart_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	// SQLの準備
	//UID,Item_ID,Quantity
	ins, err := db.Prepare("INSERT INTO cartlist (Cart_ID,Session_Key)VALUES(?,?)")
	if err != nil {
		return err
	}
	defer ins.Close()
	// SQLの実行
	_, err = ins.Exec(Cart_ID, SessionKey)
	if err != nil {
		return err
	}
	return nil
}

func Get_Cart_ID(OldSessionKey string) (Cart_ID string, err error) {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	rows, err := db.Query("SELECT CartId FROM cartlist WHERE SessionKey = ?", OldSessionKey)
	if err != nil {
		return "none", err
	}

	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Cart_ID)

		if err != nil {
			return "none", err
		}
	}
	return Cart_ID, nil
}

func (c Cart_Request_Payload) Update_Cart(Cart_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare("UPDATE cart SET Quantity = ? WHERE CartId = ? AND Item_ID = ?")
	if err != nil {
		return err
	}
	// SQLの実行
	_, err = ins.Exec(c.Quantity, Cart_ID, c.Item_ID)
	if err != nil {
		return err
	}
	defer ins.Close()
	return nil
}

func Delete_Cart(Cart_ID string, Item_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()

	// SQLの実行
	ins, err := db.Prepare("DELETE FROM cart WHERE Cart_ID = ? AND Item_ID = ?")
	if err != nil {
		return err
	}
	// SQLの実行
	_, err = ins.Exec(Cart_ID, Item_ID)
	if err != nil {
		return err
	}
	defer ins.Close()
	return nil
}

func Delete_Item_From_Cart(Item_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare("DELETE FROM cart WHERE Item_ID = ?")
	if err != nil {
		return err
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(Item_ID)
	if err != nil {
		return err
	}
	return nil
}
