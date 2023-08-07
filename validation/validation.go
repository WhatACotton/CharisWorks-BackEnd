package validation

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type User struct {
	Userdata auth.UserRecord
}
type SessionManager struct {
	Session *scs.SessionManager
}

func (user *User) Verify(c *gin.Context, uid string) (authorized bool) {
	authorized = false
	// Firebaseアプリを初期化する
	conf := &firebase.Config{
		ProjectID: "iris-test-52dcd",
	}
	opt := option.WithCredentialsFile("application_default_credentials.json")
	app, err := firebase.NewApp(context.Background(), conf, opt)
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
	authorized = true
	return authorized
}

// ベーシック認証周りの設定
func Basic(r *gin.Engine) (routergroup *gin.RouterGroup) {
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		os.Getenv("AUTH_USER"): os.Getenv("AUTH_PASS"),
	}))
	return authorized
}

// CORSの設定
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

// SessionKeyの発行
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
func GetUUID() string {
	uuidObj, _ := uuid.NewUUID()
	return uuidObj.String()
}

func SessionConfig(r *gin.Engine) {
	store := cookie.NewStore([]byte(GenerateRandomKey()))
	r.Use(sessions.Sessions("mySession", store))
}

func SessionStart(c *gin.Context) (OldSessionKey string, NewSessionKey string) {
	session := sessions.DefaultMany(c, "SessionKey")
	if session.Get("SessionKey") == nil {
		SessionKey := GetUUID()
		session.Set("SessionKey", SessionKey)
		err := session.Save()
		if err != nil {
			log.Fatal(err)
		}
		return "new", SessionKey
	} else {
		SessionKey := session.Get("SessionKey")
		NewSessionKey := GetUUID()
		session.Set("SessionKey", NewSessionKey)
		session.Save()
		return SessionKey.(string), NewSessionKey
	}
}

func CartConfig(r *gin.Engine) {
	store := cookie.NewStore([]byte(GenerateRandomKey()))
	cookies := []string{"CartSessionKey", "SessionKey"}
	r.Use(sessions.SessionsMany(cookies, store))
}

func LogInStatus(c *gin.Context) bool {
	session := sessions.DefaultMany(c, "SessionKey")
	if session.Get("SessionKey") == nil {
		return false
	} else {
		return true
	}
}

func CartSessionStart(c *gin.Context) (OldSessionKey string, NewSessionKey string) {
	session := sessions.DefaultMany(c, "CartSessionKey")
	if session.Get("CartSessionKey") == nil {
		SessionKey := GetUUID()
		session.Set("CartSessionKey", SessionKey)
		err := session.Save()
		if err != nil {
			log.Fatal(err)
		}
		return "new", SessionKey
	} else {
		SessionKey := session.Get("CartSessionKey")
		NewSessionKey := GetUUID()
		session.Set("CartSessionKey", NewSessionKey)
		session.Save()
		return SessionKey.(string), NewSessionKey
	}
}
