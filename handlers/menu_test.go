package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// Helper function to create a router during testing
func getRouter(withTemplates bool) *gin.Engine {
	r := gin.Default()
	if withTemplates {
		r.LoadHTMLGlob("../templates/*")
	}
	return r
}

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	w := httptest.NewRecorder() // Create a response recorder
	r.ServeHTTP(w, req)         // Process request
	if !f(w) {
		t.Fail()
	}
}

// GET / tests
func TestMenu(t *testing.T) {

	// Setup
	r := getRouter(true)
	budgetCodes := map[string]*BudgetCode{}
	r.GET("/", Menu("", budgetCodes))

	// Test menu loads
	req, _ := http.NewRequest("GET", "/", nil)
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK // Status code check
		p, err := ioutil.ReadAll(w.Body)    // Check page content
		pageOK := err == nil && strings.Index(string(p), "<title>Upload a File</title>") > 0
		return statusOK && pageOK
	})

}
