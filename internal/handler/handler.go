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
	//registory := c.Request.Context()
	session := c.Request.Cookies()
	if len(session) == 0 {
		SessionID := database.GetUUID()
		validation.Generate(c.Writer, c.Request, SessionID)
		funcs.StoreSession(SessionID)
	}
}

func GetsessionId(c *gin.Context) string {
	//registory := c.Request.Context()
	session := c.Request.Cookies()
	if len(session) == 0 {
		SessionID := database.GetUUID()
		validation.Generate(c.Writer, c.Request, SessionID)
		funcs.StoreSession(SessionID)
		return SessionID
	} else {
		SessionID := validation.GetSessionId(c.Writer, c.Request)
		return SessionID
	}
}
