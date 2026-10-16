package handlers

import (
	"net/http"

	"hidden-coast/internal/registry"

	"github.com/gin-gonic/gin"
)

func RegisterSMEHandlers(r *gin.Engine) {
	smeGroup := r.Group("/api/sme")
	{
		smeGroup.GET("", listSMEs)
		smeGroup.GET("/:id", getSME)
		smeGroup.GET("/:id/inventory", getSMEInventory)
	}
}

func listSMEs(c *gin.Context) {
	reg := registry.GetRegistry()
	category := c.Query("category")
	area := c.Query("area")

	if category != "" || area != "" {
		c.JSON(http.StatusOK, reg.SearchByCategory(category, area))
		return
	}

	c.JSON(http.StatusOK, reg.ListAll())
}

func getSME(c *gin.Context) {
	id := c.Param("id")
	reg := registry.GetRegistry()
	sme, ok := reg.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "SME not found"})
		return
	}
	c.JSON(http.StatusOK, sme)
}

func getSMEInventory(c *gin.Context) {
	id := c.Param("id")
	reg := registry.GetRegistry()
	sme, ok := reg.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "SME not found"})
		return
	}
	c.JSON(http.StatusOK, sme.Products)
}
