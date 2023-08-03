package funcs

import (
	"net/http"
	"unify/internal/database"
	"unify/internal/models"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func PostCartLoggedIn(usr validation.User, c *gin.Context) {
	newCartReq := new(models.CartRequestPayload)
	if err := c.BindJSON(&newCartReq); err != nil {
		return
	}
	c.JSON(http.StatusOK, database.PostCart(*newCartReq, usr.Userdata.UID))
}

func PostCartWithSession(c *gin.Context, SessionId string) {
	newCartReq := new(models.CartRequestPayload)
	if err := c.BindJSON(&newCartReq); err != nil {
		return
	}
	c.JSON(http.StatusOK, database.PostCart(*newCartReq, SessionId))
}
func GetCartWithSession(c *gin.Context, SessionId string) {
	c.JSON(http.StatusOK, database.GetCart(SessionId))
}
