package handler

import (
	"log"

	"github.com/WhatACotton/go-backend-test/internal/database"
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
func Category(c *gin.Context) {
	Category := c.Param("category")
	ItemList, err := database.ItemGetCategory(Category)
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, ItemList)
}
func Color(c *gin.Context) {
	Color := c.Param("color")
	ItemList, err := database.ItemGetColor(Color)
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, ItemList)
}
func ItemDetails(c *gin.Context) {
	ItemID := c.Query("ItemID")
	Item := new(database.Item)
	Item.ItemGet(ItemID)
	if Item.ItemID != "" {
		c.JSON(200, gin.H{"Item": Item})
	} else {
		c.JSON(404, "{Not Found}")
	}
}
func ItemMakerGet(c *gin.Context) {
	MakerName := c.Param("MakerName")
	log.Print("MakerName:", MakerName)
	Items := database.ItemGetMaker(MakerName)
	c.JSON(200, Items)

}
