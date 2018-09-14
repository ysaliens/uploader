package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ysaliens/uploader/handlers"
)

func main() {
	address := "localhost:8080"
	budgetCodesCSV := "./files/config/" + "budget_codes.csv"

	// Read Budget Codes, Category, Opex, Labels from CSV
	budgetCodes, err := handlers.CreateBudgetCodes(budgetCodesCSV)
	if err != nil {
		log.Printf("ERROR Reading budget codes: %v", err)
		return //Don't process anything if we don't know budget codes
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handlers.Menu(address, budgetCodes))
	router.POST("/download", handlers.DownloadData(address))
	router.POST("/upload", handlers.UploadFile(address, budgetCodes))

	router.Run(address)
}
