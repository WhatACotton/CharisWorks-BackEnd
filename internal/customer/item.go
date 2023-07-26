package customer

import (
	"net/http"
	"unify/internal/handler"
	"unify/internal/models"

	"github.com/gin-gonic/gin"
)

func PostItem(c *gin.Context) {
	var newItem models.Item
	if err := c.BindJSON(&newItem); err != nil {
		return
	}
	c.JSON(http.StatusOK, handler.PostItem(newItem))
}

func PatchItem(c *gin.Context) {
	var newReq models.PatchRequestPayload
	if err := c.BindJSON(&newReq); err != nil {
		return
	}
	requesttype := newReq.Attribute
	if requesttype == "Name" ||
		requesttype == "Description" ||
		requesttype == "Keyword" {
		newReq.Isint = false
	} else {
		newReq.Isint = true
	}
	c.JSON(http.StatusOK, handler.PatchItem(newReq))
}

func GetItemList(c *gin.Context) {
	c.JSON(http.StatusOK, handler.GetItemList())
}

func GetItem(c *gin.Context) {
	id := c.Query("id")
	c.JSON(http.StatusOK, handler.GetItem(id))
}

func DeleteItem(c *gin.Context) {
	id := c.Query("id")
	c.JSON(http.StatusOK, handler.DeleteItem(id))
}
