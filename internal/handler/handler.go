package handler

import (
	"net/http"
	"unify/internal/auth"
	"unify/internal/database"
	"unify/internal/funcs"
	"unify/validation"

	"github.com/gin-gonic/gin"
)

func GetItem(c *gin.Context) {
	id := c.Query("id")
	funcs.GetItem(c, id)
}

func GetItemList(c *gin.Context) {
	funcs.GetItemList(c)
}

func Transaction(c *gin.Context) {
	requestMethod := http.MethodGet
	switch request := requestMethod; request {
	case "GET":
		auth.GetTransaction(c)
	case "POST":
		auth.PostTransaction(c)
	}
}
func SessionStart(c *gin.Context) {
	registory := c.Request.Context()
	if registory.Value(registory) == nil {
		sessionId := database.GetUUID()
		validation.Generate(c.Writer, c.Request, sessionId)
		//funcs.StoreSession(sessionId)
	} else {
		//sessionID := validation.GetSessionId(c.Writer, c.Request)
		//ExpiredDate := funcs.Getsession(sessionID)
		//c.JSON(http.StatusOK, sessionID)

	}
}
