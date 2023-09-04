package validation

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type UserReqPayload struct {
	Email   string `json:"contact"`
	UID     string `json:"uid"`
	IdToken string `json:"idToken"`
}

func (user *UserReqPayload) Verify(c *gin.Context) bool {
	// Firebaseアプリを初期化する
	conf := &firebase.Config{
		ProjectID: "iris-test-52dcd",
	}
	opt := option.WithCredentialsFile("application_default_credentials.json")
	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	Token := verifyIDToken(c, app, user.IdToken)
	user.Email = Token.Claims["email"].(string)
	user.UID = Token.Claims["user_id"].(string)
	log.Printf("Successfully get \nemail: %v\nUID: %v\n", user.Email, user.UID)
	return true
}
func verifyIDToken(ctx context.Context, app *firebase.App, idToken string) *auth.Token {
	// [START verify_id_token_golang]
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	// [END verify_id_token_golang]

	return token
}
func (user *UserReqPayload) DeleteCustomer(c *gin.Context) {
	ctx := c.Request.Context()
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

	err = client.DeleteUser(ctx, user.UID)
	if err != nil {
		log.Fatalf("error deleting user: %v\n", err)
	}
	log.Printf("Successfully deleted user: %s\n", user)

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

func SessionEnd(c *gin.Context) (OldSessionKey string) {
	session := sessions.DefaultMany(c, "SessionKey")
	OldSessionKey = session.Get("SessionKey").(string)
	session.Clear()
	session.Save()
	return OldSessionKey
}

func CartSessionEnd(c *gin.Context) (OldSessionKey string) {
	session := sessions.DefaultMany(c, "CartSessionKey")
	OldSessionKey = session.Get("CartSessionKey").(string)
	session.Clear()
	session.Save()
	return OldSessionKey
}

func SessionConfig(r *gin.Engine) {
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
func Get_Cart_Session(c *gin.Context) (Cart_Session_Key string) {
	session := sessions.DefaultMany(c, "CartSessionKey")
	if session.Get("CartSessionKey") == nil {
		return "new"
	} else {
		Cart_Session_Key = session.Get("CartSessionKey").(string)
		return Cart_Session_Key
	}
}

func Set_Cart_Session(c *gin.Context, Cart_Session_Key string) {
	session := sessions.DefaultMany(c, "CartSessionKey")
	session.Set("CartSessionKey", Cart_Session_Key)
	session.Save()
}
