package database

import "github.com/pkg/errors"

type Cart struct {
	//UUID
	Cart_ID string `json:"Cart_ID"`
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
	rows, err := db.Query("SELECT Order,Item_ID,Quantity,Info_ID,Status,Price,Name,Stock FROM Cart_List WHERE Cart_ID = ?", Cart_ID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting prepare Cart_ID")
	}
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		err := rows.Scan(&Cart.Order, &Cart.Item_ID, &Cart.Quantity, &Cart.Info_ID, &Cart.Status, &Cart.Price, &Cart.Item_Name, &Cart.Stock)
		if err != nil {
			return nil, errors.Wrap(err, "error in scanning Cart_ID")
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
		return nil, errors.Wrap(err, "error in getting Cart_ID")
	}
	for _, _Cart := range Carts {
		err := _Cart.Get_Cart()
		if err != nil {
			return nil, errors.Wrap(err, "error in getting Cart_ID")
		}
		Carts = append(Carts, _Cart)
	}
	return Carts, nil
}

func (c *Cart) Get_Cart() error {
	err := c.Get_Item_Info()
	if err != nil {
		return err
	}
	Item_List := new(Item_List)
	err = Item_List.Get_Item_List(c.Item_ID)
	if err != nil {
		return err
	}
	c.Status = Item_List.Status
	return nil
}
func (c *Cart) Get_Item_Info() error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, err := db.Query("SELECT Item_ID,Info_ID,Status,Name,Stock FROM Item_List WHERE Item_ID = ?", c.Item_ID)
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

func (c Cart_Request_Payload) Cart(Cart_ID string) error {
	Item_List := new(Item_List)
	err := Item_List.Get_Item_List(c.Item_ID)
	if err != nil {
		return err
	}
	//リクエストされた商品が登録可能か判定
	if Item_List.Status == "Available" {
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

				if Item_List.Stock >= c.Quantity {
					err := c.Update_Cart(Cart_ID)
					if err != nil {
						return err
					}
				} else {
					return errors.New("stock is not enough")
				}
			}
		} else {
			err := c.Post_Cart(Cart_ID)
			if err != nil {
				return err
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
