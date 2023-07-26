package database

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func GetDate() string {
	const template = "2006-01-02 15:04:05"
	t := time.Now()
	t.Format(template)
	return t.String()
}
func GettransactionId() string {
	uuidObj, _ := uuid.NewUUID()
	return uuidObj.String()
}

type PatchRequestPayload struct {
	ID        string `json:"id"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

func (patchItem PatchRequestPayload) Patch(table string, isint bool, where string) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

	// SQLの準備
	upd, err := db.Prepare("UPDATE ? SET ? = ? WHERE ? = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer upd.Close()

	if isint == true {
		value, err := strconv.Atoi(patchItem.Value)
		if err != nil {
			log.Fatal(err)
		}
		// SQLの実行
		_, err = upd.Exec(table, patchItem.Attribute, value, where, patchItem.ID)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// SQLの実行
		_, err = upd.Exec(table, patchItem.Attribute, patchItem.Attribute, where, patchItem.ID)
		if err != nil {
			log.Fatal(err)
		}
	}

}
func Delete(table string, where string, id string) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")
	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

	// SQLの実行
	del, err := db.Prepare("DELETE FROM ? WHERE ? = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer del.Close()

	// SQLの実行
	_, err = del.Exec(table, where, id)
	if err != nil {
		log.Fatal(err)
	}
}
