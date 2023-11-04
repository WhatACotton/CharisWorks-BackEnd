package validation

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type CustomerReqPayload struct {
	Email         string
	UserID        string
	EmailVerified bool
}

// HeaderのAuthenticationに入っているJWTからEmail,UserID,EmailVerifiedを取得
func (user *CustomerReqPayload) VerifyCustomer(c *gin.Context) bool {
	// Firebaseアプリを初期化する
	conf := &firebase.Config{
		ProjectID: "charisworks-a1ef5",
	}
	opt := option.WithCredentialsFile("application_default_credentials.json")
	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	IdToken := get_IdToken(c.Request)
	Token := verifyIDToken(c, app, IdToken)
	user.Email = Token.Claims["email"].(string)
	user.UserID = Token.Claims["user_id"].(string)
	user.EmailVerified = Token.Claims["email_verified"].(bool)
	log.Printf("Successfully get \nemail: %v\nUserID: %v\nEmail_Verified: %v\n", user.Email, user.UserID, user.EmailVerified)
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

func get_IdToken(r *http.Request) (token string) {
	// Authorizationヘッダーの値を取得
	token = r.Header.Get("Authorization")
	// Bearerトークンの抽出
	return token
	// tokenを利用して任意の処理を行う
}
