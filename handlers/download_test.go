package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// GET /download tests
func TestDownloadData(t *testing.T) {

	// Setup
	r := getRouter(true)
	r.GET("/download", DownloadData(""))

	// @TODO Mock db entries with actual data for bottom test cases

	// Test vessel name requirement
	form := url.Values{}
	req, _ := http.NewRequest("GET", "/download", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = form
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK // Status code check
		p, err := ioutil.ReadAll(w.Body)    // Check page content
		pageOK := err == nil && strings.Index(string(p), "Vessel Name is required.") > 0
		return statusOK && pageOK
	})

	// Test year requirement
	form = url.Values{}
	form.Add("name", "VESSELNAME")
	form.Add("year", "-5")
	req, _ = http.NewRequest("GET", "/download", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = form
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "Year must be positive.") > 0
		return statusOK && pageOK
	})

	// Test all data by vessel name
	form = url.Values{}
	form.Add("name", "VESSELNAME")
	req, _ = http.NewRequest("GET", "/download", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = form
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "No records found") > 0
		return statusOK && pageOK
	})

	// Test all data by vessel name+year
	form = url.Values{}
	form.Add("name", "VESSELNAME")
	form.Add("year", "5")
	req, _ = http.NewRequest("GET", "/download", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = form
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "No records found") > 0
		return statusOK && pageOK
	})

	// Test all data by vessel name+year+code
	form = url.Values{}
	form.Add("name", "VESSELNAME")
	form.Add("year", "2015")
	form.Add("code", "1111")
	req, _ = http.NewRequest("GET", "/download", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = form
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "No records found") > 0
		return statusOK && pageOK
	})

}
