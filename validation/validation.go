package validation

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type User struct {
	Userdata auth.UserRecord
}

func (user User) Verify(c *gin.Context) (authorized bool) {
	authorized = false
	// Firebaseアプリを初期化する
	conf := &firebase.Config{
		ProjectID: "iris-test-52dcd",
	}

	opt := option.WithCredentialsFile("application_default_credentials.json")
	app, err := firebase.NewApp(context.Background(), conf, opt)
	uid := c.Query("uid")
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// Get an auth client from the firebase.App
	client, err := app.Auth(c)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	u, err := client.GetUser(c, uid)
	if err != nil {
		log.Fatalf("error getting user %s: %v\n", uid, err)
	}
	log.Printf("Successfully fetched user data: %v\n", u)
	user.Userdata = *u
	c.JSON(http.StatusOK, *&u.UID)
	authorized = true
	return authorized
}

func Basic(r *gin.Engine) (routergroup *gin.RouterGroup) {
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		os.Getenv("AUTH_USER"): os.Getenv("AUTH_PASS"),
	}))
	return authorized
}

func CORS(r *gin.Engine) {
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
			"Access-Control-Allow-Headers",
			"Authorization",
			"Access-Control-Allow-Credentials",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: false,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	}))
}
