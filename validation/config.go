package validation

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CORSの設定
func CORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		// アクセス許可するオリジン
		AllowOrigins: []string{
			os.Getenv("FRONTEND_URL"),
		},
		// アクセス許可するHTTPメソッド
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"PATCH",
			"DELETE",
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
	log.Print(r)
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
	r.Use(sessions.Sessions("SessionKey", store))
}

// アカウント関連のロギング
func LoginLogging(Text string) {
	file, err := os.OpenFile("accountlog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime)
	logger.SetOutput(file)
	logger.Println("UserID :", Text)
}

// 送金関連のロギング
func TransferLogging(Text string) {
	file, err := os.OpenFile("transferlog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime)
	logger.SetOutput(file)
	logger.Println("UserID :", Text)
}
