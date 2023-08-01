package validation

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
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
		AllowCredentials: true,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	}))
}

// Note: Don't store your key in your source code. Pass it via an
// environmental variable, or flag (or both), and don't accidentally commit it
// alongside your code. Ensure your key is sufficiently random - i.e. use Go's
// crypto/rand or securecookie.GenerateRandomKey(32) and persist the result.
// Ensure SESSION_KEY exists in the environment, or sessions will fail.
var store = sessions.NewCookieStore([]byte(GenerateRandomKey()))

func Generate(w http.ResponseWriter, r *http.Request, uuid string) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set some session values.
	session.Values["sessionID"] = uuid

	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func GetSessionId(w http.ResponseWriter, r *http.Request) (sessionId string) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessionId = session.Values["sessionID"].(string)
	return sessionId
}

func GenerateRandomKey() (sessionKey string) {
	// 32バイトのランダムなバイト列を生成する
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	// バイト列をBase64エンコードして、文字列に変換する
	sessionKey = base64.URLEncoding.EncodeToString(key)
	return
}
