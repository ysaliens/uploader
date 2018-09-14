package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ysaliens/uploader/models"
)

// DataForm - Get Data Form
type DataForm struct {
	Name       string `form:"name" binding:"required"`
	Year       int    `form:"year"`
	BudgetCode string `form:"code"`
}

// DownloadData - GET /download - Query and download ship data as JSON
func DownloadData(address string) gin.HandlerFunc {
	return func(c *gin.Context) {
		upload := "http://" + address + "/upload"
		download := "http://" + address + "/download"

		// Bind and verify user input from page
		var form DataForm
		if c.ShouldBind(&form) != nil {
			log.Printf("Bad User Input - Name")
			c.HTML(http.StatusOK, "menu.tmpl", gin.H{"upload": upload, "download": download, "status": "Vessel Name is required."})
			return
		}
		if form.Year < 0 {
			log.Printf("Bad User Input - Year")
			c.HTML(http.StatusOK, "menu.tmpl", gin.H{"upload": upload, "download": download, "status": "Year must be positive."})
			return
		}
		//log.Printf("form.Name: %s form.Year: %d form.BudgetCode: %s", form.Name, form.Year, form.BudgetCode)

		// Search database for records
		var e error
		var records []models.Record
		manager := models.MongoDBConnection{}
		if len(form.BudgetCode) == 0 && form.Year == 0 {
			records, e = manager.FetchByName(form.Name)
		} else if len(form.BudgetCode) == 0 && form.Year > 0 {
			records, e = manager.FetchByNameYear(form.Name, strconv.Itoa(form.Year))
		} else if len(form.BudgetCode) > 0 && form.Year > 0 {
			records, e = manager.FetchByNameYearCode(form.Name, strconv.Itoa(form.Year), form.BudgetCode)
		}

		// Return a JSON list of records
		if e != nil || len(records) == 0 {
			c.HTML(http.StatusOK, "menu.tmpl", gin.H{"upload": upload, "download": download, "status": "No records found"})
			return
		}
		c.JSON(http.StatusOK, records)
	}
}
