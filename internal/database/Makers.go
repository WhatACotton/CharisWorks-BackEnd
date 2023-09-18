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

func (m Maker) MakerItemDetailsModyfy(i ItemDetails) {
	//MakerItemDetailsの修正
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
func MakerItemDetailsCreate() {

}
func ItemCreate() {
	// Itemの作成
}
func ItemChange() {
	//引数は変更前,変更後
	//紐付いているItemの変更
}
func ItemChangeStatus() {
	// Itemの状態の変更
}
