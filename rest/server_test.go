package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	database "github.com/piotrostr/realtime/db"
	"github.com/piotrostr/realtime/rest"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, 201, w.Code)
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

	req, err := http.NewRequest("DELETE", "/delete", nil)
	assert.Nil(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}
