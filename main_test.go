package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test case for a simple success case
func TestSuccessfulResponse(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/api/order?quantity=123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Edge cases

// Test case for minimizing number of packs shipped
// In this example, 2x250 or 1x500 works, but we want to prefer 1x500.
func TestEdgeCase1(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/api/order?quantity=251&sizes=250,500,1000,2000,5000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var m map[string]int
	json.Unmarshal(w.Body.Bytes(), &m)
	assert.Equal(t, m["500"], 1)
}

// Test case for minimizing excess items shipped.
// In this example, 1x1000 or 3x250 works, but we want to send (1x500, 1x250)
// so that the excess items are minimum (in this case, 249 would be minimum).
// Also, comparing 3x250 and (1x500, 1x250), we chose the latter for minimizing
// number of packs shipped.
func TestEdgeCase2(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/api/order?quantity=501&sizes=250,500,1000,2000,5000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var m map[string]int
	json.Unmarshal(w.Body.Bytes(), &m)
	assert.Equal(t, m["500"], 1)
	assert.Equal(t, m["250"], 1)
}

// Another test case for minimizing excess items sent.
// In this case, 1x2000 and (3x600, 1x300) works.
// We want to choose 1x2000 because the latter has more excess items (2100 in total).
func TestEdgeCase3(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/api/order?quantity=1801&sizes=300,600,2000,5000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var m map[string]int
	json.Unmarshal(w.Body.Bytes(), &m)
	assert.Equal(t, m["2000"], 1)
}

// Test case for input where no order quantity is specified
func TestNoQuantityInput(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/api/order", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test case for input where pack sizes are badly specified (e.g. has a letter)
func TestBadSizesInput(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/api/order?quantity=1&sizes=1,2,3a", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test case for input where order quantity is negative
func TestBadQuantityInput(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/api/order?quantity=-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
