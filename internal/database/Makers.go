package database

type Maker struct {
	MadeBy          string `json:"MadeBy"`
	Description     string `json:"Description"`
	StripeAccountID string `json:"StripeAccountID"`
}

func MakerAccountCreate(m Maker) {
	//stripeAccountを元にアカウントを作成する。
	//アカウントの画像はpathで指定して特定の場所に保存しておく必要がある。
	db := ConnectSQL()
	ins, _ := db.Prepare(`
	INSERT
	INTO
		MakersDetails
		(MadeBy,Description,StripeAccountID)
	VALUES
		(?,?,?)`)
	ins.Exec(m.MadeBy, m.Description, m.StripeAccountID)
	defer ins.Close()
}

func MakerAccountGet(MadeBy string) {
	//stripeAccountID以外を取得してくる。
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	rows, _ := db.Query(`
	SELECT
		MadeBy,
		Description
	FROM
		MakersDetails
	WHERE
		MadeBy = ?`, MadeBy)
	defer rows.Close()
	// SQLの実行
	for rows.Next() {
		rows.Scan(&MadeBy)
	}
}
func (m Maker) MakerAccountModyfy() {
	//アカウント情報の修正
	db := ConnectSQL()
	defer db.Close()
	// SQLの実行
	db.Exec(`
	UPDATE
		MakersDetails
	SET
		MadeBy = ?,
		Description = ?
	WHERE
		StripeAccountID = ?`, m.MadeBy, m.Description, m.StripeAccountID)
}

// 削除は取引済みの商品の対応・アカウントなどの対応が必要なためできない。(すべてが終わったあとに実装するかもしれない。)
func MakerAccountDelete() {
	//アカウント削除
	//関連しているテーブルに影響を与えざるを得ないので、あまり使いたくない。
}
func MakerGetStripeID(Maker MadeBy) (StripeAccountID string) {
	db := ConnectSQL()
	defer db.Close()
	rows, _ := db.Query(`
	SELECT 
		StripeAccountID
	FROM 
		MakersDetails
	WHERE 
		MadeBy = ?`, Maker)
	rows.Scan(&StripeAccountID)
	defer rows.Close()
	return StripeAccountID
}
func (m MadeBy) MakerItemsGet() (Items Items) {
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
			MadeBy = ?`,
		m)
	defer db.Close()
	for rows.Next() {
		Item := new(Item)
		rows.Scan(&Item.ItemID, &Item.Status, &Item.ItemName, &Item.Price, &Item.Stock, &Item.ItemOrder)
		Items = append(Items, *Item)
	}
	return Items

}
func (m MadeBy) MakerItemCreate(Item Item) {
	db := ConnectSQL()
	//UID,ItemID,Quantity
	ins, _ := db.Prepare(`
	INSERT 
	INTO 
		Item 
		(DetailsID , Status , Price , Stock , ItemName , MadeBy) 
		VALUES 
		(? , ? , ? , ? , ? , ?)`)
	defer ins.Close()
	// SQLの実行
	ins.Exec(Item.DetailsID, Item.Status, Item.Price, Item.Stock, Item.ItemName, Item.MadeBy)
}
