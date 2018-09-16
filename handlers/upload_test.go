package handlers

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// POST /upload tests
func TestUpload(t *testing.T) {

	// Setup
	r := getRouter(true)
	budgetCodes := map[string]*BudgetCode{}
	r.POST("/upload", UploadFile("", budgetCodes))

	// Errors on bad file
	form := url.Values{}
	form.Add("uploadfile", "File1.txt")
	req, _ := http.NewRequest("POST", "/upload", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "multipart/form-data")
	req.PostForm = form
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK // Status code check
		p, err := ioutil.ReadAll(w.Body)    // Check page content
		pageOK := err == nil && strings.Index(string(p), "ERROR: File upload failed.") > 0
		return statusOK && pageOK
	})

	// Errors on bad extension
	testFile := "../files/test/WrongExtension.txt"
	file, err := os.Open(testFile)
	if err != nil {
		t.Errorf("%s: %s", err, testFile)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("uploadfile", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()
	req, _ = http.NewRequest("POST", "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK // Status code check
		p, err := ioutil.ReadAll(w.Body)    // Check page content
		pageOK := err == nil && strings.Index(string(p), "ERROR: Unsupported file extension. Please use .xlsx or .zip") > 0
		return statusOK && pageOK
	})

	// Skips .xlsx files with bad format
	// @ TODO - Check nothing got added in mock db
	testFile = "../files/test/BadFormat.xlsx"
	file, err = os.Open(testFile)
	if err != nil {
		t.Errorf("%s: %s", err, testFile)
	}
	defer file.Close()
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	part, _ = writer.CreateFormFile("uploadfile", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()
	req, _ = http.NewRequest("POST", "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK // Status code check
		p, err := ioutil.ReadAll(w.Body)    // Check page content
		pageOK := err == nil && strings.Index(string(p), "Upload Successful") > 0
		return statusOK && pageOK
	})

	// Correctly processes an xlsx file
	// @ TODO - Check mock db has records
	testFile = "../files/test/test.xlsx"
	file, err = os.Open(testFile)
	if err != nil {
		t.Errorf("%s: %s", err, testFile)
	}
	defer file.Close()
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	part, _ = writer.CreateFormFile("uploadfile", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()
	req, _ = http.NewRequest("POST", "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK // Status code check
		p, err := ioutil.ReadAll(w.Body)    // Check page content
		//t.Errorf("%s", string(p))
		pageOK := err == nil && strings.Index(string(p), "Upload Successful") > 0
		return statusOK && pageOK
	})

}
