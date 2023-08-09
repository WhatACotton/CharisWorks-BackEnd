package database

import "unify/internal/models"

func PostTransactionList(CartId string, UID string, TransactionDate string) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare("INSERT INTO transactionlist VALUES(?,?,?)")
	if err != nil {
		//log.Fatal(err)
		panic(err.Error())

	}
	defer ins.Close()

	// SQLの実行
	_, err = ins.Exec(
		CartId,
		UID,
		TransactionDate,
	)
	if err != nil {
		//log.Fatal(err)
		panic(err.Error())
	}
}

func PostTransaction(Carts []models.Cart) {
	// データベースのハンドルを取得する
	db := ConnectSQL()
	defer db.Close()

	// SQLの準備
	ins, err := db.Prepare("INSERT INTO transactionlist VALUES(?,?,?)")
	if err != nil {
		//log.Fatal(err)
		panic(err.Error())

	}
	defer ins.Close()
	for _, Cart := range Carts {
		// SQLの実行
		_, err = ins.Exec(
			Cart.CartId,
			Cart.InfoId,
			Cart.Quantity,
		)
		if err != nil {
			//log.Fatal(err)
			panic(err.Error())

		}
	}
}
