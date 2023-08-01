package funcs

import (
	"net/http"
	"unify/internal/database"
	"unify/internal/models"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func PostCartLoggedIn(usr validation.User, c *gin.Context) {
	var newCartReq models.CartRequestPayload
	if err := c.BindJSON(&newCartReq); err != nil {
		return
	}
	c.JSON(http.StatusOK, database.PostCart(newCartReq, usr.Userdata.UID))
}
