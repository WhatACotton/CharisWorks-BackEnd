package handler

import (
	"log"
	"unify/internal/database"

	"github.com/gin-gonic/gin"
)

func Top(c *gin.Context) {
	TopItemList, err := database.ItemGetTop()
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, TopItemList)
}
func ALL(c *gin.Context) {
	ItemList, err := database.ItemGetALL()
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, ItemList)
}
func ItemDetails(c *gin.Context) {
	ItemID := c.Query("ItemID")
	InfoID, err := database.ItemDetailsIDGet(ItemID)
	if err != nil {
		log.Print(err)
	}
	if InfoID != "Couldn't get" {
		ItemDetails, err := database.ItemDetailsGet(InfoID)
		if err != nil {
			log.Print(err)
		}
		c.JSON(200, ItemDetails)
	} else {
		c.JSON(404, "{Not Found}")
	}
}
func Category(c *gin.Context) {
	Category := c.Param("category")
	ItemList, err := database.ItemCategoryGet(Category)
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, ItemList)
}
func Color(c *gin.Context) {
	Color := c.Param("color")
	ItemList, err := database.ItemColorGet(Color)
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, ItemList)
}
