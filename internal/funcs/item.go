package funcs

import (
	"net/http"
	"unify/internal/database"

	"github.com/gin-gonic/gin"
)

func GetItemList(c *gin.Context) { c.JSON(http.StatusOK, database.GetItemList()) }

func GetItem(c *gin.Context, id string) { c.JSON(http.StatusOK, database.GetItem(id)) }
