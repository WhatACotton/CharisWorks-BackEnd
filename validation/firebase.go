package validation

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type UserReqPayload struct {
	Email          string `json:"contact"`
	UID            string `json:"uid"`
	IdToken        string `json:"idToken"`
	Email_Verified bool   `json:"email_verified"`
	Cart_ID        string
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
	user.Email_Verified = Token.Claims["email_verified"].(bool)
	log.Printf("Successfully get \nemail: %v\nUID: %v\nEmail_Verified: %v\n", user.Email, user.UID, user.Email_Verified)
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
