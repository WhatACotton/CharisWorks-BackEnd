package main

import (
	"database/sql"
	"log"
	"os"
	"time"
	"unify/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		// アクセス許可するオリジン
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		// アクセス許可するHTTPメソッド
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"PATCH",
		},
		// 許可するHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Content-Type",
			"Access-Control-Allow-Origin",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: false,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	}))
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", os.Getenv("MYSQL_PASS")+":"+os.Getenv("MYSQL_USER")+"@tcp(localhost:3306)/go_test")
	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

	// 実際に接続する
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return
	} else {
		log.Println("データベース接続完了")
	}
	r.POST("/customer", handler.PostCustomer)
	r.GET("/customer", handler.GetCustomer)
	r.DELETE("/customer", handler.DeleteCustomer)
	r.PATCH("/customer", handler.UpdateCustomerCustomer)

	r.POST("/transaction", handler.PostTransaction)
	r.GET("/transaction", handler.GetTransaction)
	r.DELETE("/transaction", handler.DeleteTransaction)
	r.PATCH("/transaction", handler.UpdateTransaction)

	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		os.Getenv("AUTH_USER"): os.Getenv("AUTH_PASS"),
	}))
	authorized.GET("/hello", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		c.JSON(200, gin.H{"message": "Hello " + user})
	})
	authorized.GET("/items", handler.GetItem)
	authorized.POST("/items", handler.PostItem)
	authorized.PATCH("/items", handler.PatchItem)
	authorized.DELETE("/items", handler.DeleteItem)

	r.Run(":8081") // 0.0.0.0:8080 でサーバーを立てます。

}
