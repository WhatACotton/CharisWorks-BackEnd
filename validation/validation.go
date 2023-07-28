package validation

import (
	"log"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
)

// client := &http.Client{
// 	CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 		fmt.Printf("Redirected from %s to %s\n", via[len(via)-1].URL.String(), req.URL.String())
// 		return nil
// 	},
// }

// resp, err := client.Get("http://example.com")
// if err != nil {
// 	fmt.Println(err)
// 	return
// }

// defer resp.Body.Close()

// fmt.Println(resp.StatusCode)
// authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
// 	os.Getenv("AUTH_USER"): os.Getenv("AUTH_PASS"),
// }))
// authorized.GET("/hello", func(c *gin.Context) {
// 	user := c.MustGet(gin.AuthUserKey).(string)
// 	c.JSON(200, gin.H{"message": "Hello " + user})
// })

// authorized.GET("/items", handler.GetItem)
// authorized.POST("/items", handler.PostItem)
// authorized.PATCH("/items", handler.PatchItem)
// authorized.DELETE("/items", handler.DeleteItem)

func Verify(c *gin.Context, uid string) {
	app, err := firebase.NewApp(c, nil)
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

}
