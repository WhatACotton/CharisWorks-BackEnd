package database

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetDate() string {
	t := time.Now()
	ft := t.Format("2006-01-02 15:04:05")
	return ft
}

// []uint8型の値をtime.Time型に変換する
func ConvertBytesToTime(b []uint8) time.Time {
	str := string(b)
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(i, 0)
}

type PatchRequestPayload struct {
	ID        string `json:"id"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

func ConnectSQL() (db *sql.DB) {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_USER")+":"+os.Getenv("MYSQL_PASS")+"@tcp(localhost:3306)/cart")

	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	return db
}

func TestSQL() {
	// データベースのハンドルを取得する
	mysql := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASS") + "@tcp(localhost:3306)/cart"
	log.Println(mysql)
	db, err := sql.Open("mysql", mysql)
	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}

	// 実際に接続する
	err = db.Ping()
	if err != nil {
		log.Println("データベースに接続できません。MySQLが起動しているか、環境変数が設定されているか確認してください。")
		log.Fatal(err)
		return
	} else {
		log.Println("データベース接続確認")
	}

}
