package database

type ItemInfo struct {
	InfoId      string `json:"infoid"`
	Price       int    `json:"price"`
	Name        string `json:"Name"`
	Stonesize   int    `json:"Stonesize"`
	Minlength   int    `json:"Minlength"`
	Maxlength   int    `json:"Maxlength"`
	Decsription string `json:"Description"`
	Keyword     string `json:"Keyword"`
}

func PostTransactionList(CartId string, UID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare("INSERT INTO transactionlist (Cart_ID,UID)VALUES(?,?,?)")
	if err != nil {
		return err
	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(
		CartId,
		UID,
	)
	if err != nil {
		return err
	}
	return nil
}

func PostTransaction(Carts []Cart, Cart_ID string) error {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare("INSERT INTO transactionlist VALUES(?,?,?)")
	if err != nil {
		return err

	}
	defer ins.Close()
	for _, Cart := range Carts {
		// SQLの実行
		_, err = ins.Exec(
			Cart_ID,
			Cart.Info_ID,
			Cart.Quantity,
		)
		if err != nil {
			return err

		}
	}
	return nil
}
