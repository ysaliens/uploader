package handlers

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"github.com/ysaliens/uploader/models"
)

// UploadFile - POST /upload - Upload and process to db a .xlsx or .zip
func UploadFile(address string, budgetCodes map[string]*BudgetCode) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := "Upload Successful"
		upload := "http://" + address + "/upload"
		download := "http://" + address + "/download"

		file, err := c.FormFile("uploadfile")
		if err != nil {
			log.Println(err)
			status = "ERROR: File upload failed."
			c.HTML(http.StatusOK, "menu.tmpl", gin.H{"upload": upload, "download": download, "status": status})
			return
		}

		path := "./files/temp/" + file.Filename
		ex := filepath.Ext(path)
		if ex != ".xlsx" && ex != ".zip" {
			log.Println("Unsupported file extension ", ex)
			status = "ERROR: Unsupported file extension. Please use .xlsx or .zip"
			c.HTML(http.StatusOK, "menu.tmpl", gin.H{"upload": upload, "download": download, "status": status})
			return
		}

		log.Printf("File: %v\n", file.Filename)

		// @TODO - Don't save files to disk, XLSX reader lib requires it for now
		c.SaveUploadedFile(file, path)

		// Handle ZIP files here
		if ex == ".zip" {
			zipDir, err := ioutil.TempDir("./files/temp/", "zip")
			if err != nil {
				log.Println(err)
				delete(path)
				status = "Failed to unzip file."
				c.HTML(http.StatusOK, "menu.tmpl", gin.H{"upload": upload, "download": download, "status": status})
				return
			}

			unzip(path, zipDir)
			err = filepath.Walk(zipDir, func(currentFile string, info os.FileInfo, err error) error {
				// @TODO This is already set for using go routines to insert all of these in parallel
				// This should be balanced between insertion speed and server load
				if filepath.Ext(currentFile) == ".xlsx" {
					log.Println("Processing file ", currentFile)
					processXLSX(currentFile, budgetCodes)
					delete(currentFile)
				}
				return nil
			})
			if err != nil {
				log.Println(err)
			}
			delete(zipDir)

		} else { // Not a zip, read Excel file and process data
			processXLSX(path, budgetCodes)
		}

		delete(path)
		c.HTML(http.StatusOK, "menu.tmpl", gin.H{"upload": upload, "download": download, "status": status})

	}
}

func processXLSX(file string, budgetCodes map[string]*BudgetCode) {
	// @TODO This is currently very brittle if xlsx format is off
	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		log.Println(err)
		return
	}

	results := []models.Record{}

	// Read all but last sheet (Last sheet is SSR)
	for i := 0; i < (len(xlFile.Sheets) - 1); i++ {
		sheet := xlFile.Sheets[i]

		// Headers
		re := regexp.MustCompile("(\\d{4})") //Get year from "OPEX Monthly Spend (YEAR)"
		vesselName := sheet.Rows[0].Cells[1].String()
		yearString := sheet.Rows[1].Cells[1].String()
		if len(yearString) == 0 || len(vesselName) == 0 {
			return // If we can't find name or year, sheet is bad
		}
		year := re.FindAllString(yearString, -1)[0]
		log.Printf("Name: %s Year: %s\n", vesselName, year)

		// Body starts on row 7 in all files examined, uniformity assumption
		for i := 6; i < len(sheet.Rows); i++ {
			row := sheet.Rows[i]

			// Only take complete data, ignore category headers & empty lines
			// If first 2 months don't have data, it's a category or empty line
			if row.Cells[2].String() == "" && row.Cells[3].String() == "" {
				continue
			}

			// Extract body data
			regExpBudgetCode := regexp.MustCompile("[0-9.]+")
			budgetCodeList := regExpBudgetCode.FindAllString(row.Cells[1].String(), -1)
			budgetCode := ""
			if len(budgetCodeList) > 0 {
				budgetCode = budgetCodeList[0]
			} else {
				continue // This row had no budget code, must be a category/total, skip it
			}

			// Create model record object
			r := models.Record{}
			r.Name = vesselName
			r.Year = year
			r.Opex = budgetCodes[budgetCode].Opex
			r.Category = budgetCodes[budgetCode].Category
			r.BudgetCode = budgetCode
			r.BudgetDesc = budgetCodes[budgetCode].Label
			r.Jan = row.Cells[2].String()
			r.Feb = row.Cells[3].String()
			r.Mar = row.Cells[4].String()
			r.Apr = row.Cells[5].String()
			r.May = row.Cells[6].String()
			r.Jun = row.Cells[7].String()
			r.Jul = row.Cells[8].String()
			r.Aug = row.Cells[9].String()
			r.Sep = row.Cells[10].String()
			r.Oct = row.Cells[11].String()
			r.Nov = row.Cells[12].String()
			r.Dec = row.Cells[13].String()
			r.TTL = row.Cells[14].String()

			//log.Println(r.Name,r.Year,r.Opex,r.Category,r.BudgetCode,r.BudgetDesc,r.Jan,r.Feb,r.Mar,r.Apr,r.May,r.Jun,r.Jul,r.Aug,r.Sep,r.Oct,r.Nov,r.Dec,r.Ttl)
			results = append(results, r)
		}

		// Save sheet to database
		manager := models.MongoDBConnection{}
		manager.Create(results)

	}
}

func delete(file string) {
	err := os.Remove(file)
	log.Println("Deleting ", file)
	if err != nil {
		log.Println(err)
	}
}

// From http://blog.ralch.com/tutorial/golang-working-with-zip/
func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	defer reader.Close()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()
		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}
