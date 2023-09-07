package handler

import (
	"log"
	"unify/internal/database"

	"github.com/gin-gonic/gin"
)

func Top(c *gin.Context) {
	Top_Item_List, err := database.Get_Top()
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, Top_Item_List)
}
func ALL(c *gin.Context) {
	Item_List, err := database.Get_ALL()
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, Item_List)
}
func Item_Details(c *gin.Context) {
	Item_ID := c.Query("Item_ID")
	Info_ID, err := database.Get_Info_Id(Item_ID)
	if err != nil {
		log.Print(err)
	}
	if Info_ID != "Couldn't get" {
		Item_Details, err := database.Get_Item_Details(Info_ID)
		if err != nil {
			log.Print(err)
		}
		c.JSON(200, Item_Details)
	} else {
		c.JSON(404, "{Not Found}")
	}
}
func Category(c *gin.Context) {
	Category := c.Param("category")
	Item_List, err := database.Get_Item_Category(Category)
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, Item_List)
}
func Color(c *gin.Context) {
	Color := c.Param("color")
	Item_List, err := database.Get_Item_Color(Color)
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, Item_List)
}
