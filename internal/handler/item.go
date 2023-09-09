package handler

import (
	"log"
	"unify/internal/database"

	"github.com/gin-gonic/gin"
)

func Top(c *gin.Context) {
	TopItemList, err := database.GetTop()
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, TopItemList)
}
func ALL(c *gin.Context) {
	ItemList, err := database.GetALL()
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, ItemList)
}
func ItemDetails(c *gin.Context) {
	ItemID := c.Query("ItemID")
	InfoID, err := database.GetInfoId(ItemID)
	if err != nil {
		log.Print(err)
	}
	if InfoID != "Couldn't get" {
		ItemDetails, err := database.GetItemDetails(InfoID)
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
	ItemList, err := database.GetItemCategory(Category)
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, ItemList)
}
func Color(c *gin.Context) {
	Color := c.Param("color")
	ItemList, err := database.GetItemColor(Color)
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, ItemList)
}
