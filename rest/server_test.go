package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	database "github.com/piotrostr/realtime/db"
	"github.com/piotrostr/realtime/rest"
	"github.com/stretchr/testify/assert"
)

func PrintResponse(t *testing.T, req *http.Request) {
	var p []byte
	_, err := req.Body.Read(p)
	assert.Nil(t, err)
	fmt.Println(string(p))
}

func UserJsonFixture(t *testing.T) *bytes.Buffer {
	user := database.User{Name: "Piotr", Age: 22}
	userJson, err := json.Marshal(user)
	assert.Nil(t, err)
	return bytes.NewBuffer([]byte(userJson))
}

func TestCreate(t *testing.T) {
	router := rest.GetRouter()
	w := httptest.NewRecorder()

	user := UserJsonFixture(t)
	req, err := http.NewRequest("POST", "/create", user)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)

	PrintResponse(t, req)

	assert.Equal(t, 201, w.Code)
}

func TestReadOne(t *testing.T) {
	router := rest.GetRouter()
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/read/Piotr", nil)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRead(t *testing.T) {
	router := rest.GetRouter()
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/read", nil)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestUpdate(t *testing.T) {
	router := rest.GetRouter()
	w := httptest.NewRecorder()

	user := UserJsonFixture(t)
	req, err := http.NewRequest("PUT", "/update", user)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}

func TestDelete(t *testing.T) {
	router := rest.GetRouter()
	w := httptest.NewRecorder()

	// create user for delete to work
	user := UserJsonFixture(t)
	createReq, err := http.NewRequest("POST", "/create", user)
	assert.Nil(t, err)
	router.ServeHTTP(w, createReq)
	assert.Equal(t, 201, w.Code)

	w = httptest.NewRecorder()
	user = UserJsonFixture(t)
	deleteReq, err := http.NewRequest("DELETE", "/delete", user)
	assert.Nil(t, err)
	router.ServeHTTP(w, deleteReq)

	assert.Equal(t, 204, w.Code)
}
