package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Menu - GET / - Main menu page
func Menu(address string, budgetCodes map[string]*BudgetCode) gin.HandlerFunc {
	return func(c *gin.Context) {

		upload := "http://" + address + "/upload"
		download := "http://" + address + "/download"
		c.HTML(http.StatusOK, "menu.tmpl", gin.H{"upload": upload, "download": download})

	}
}
