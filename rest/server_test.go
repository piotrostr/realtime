package rest

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	req := httptest.NewRequest("POST", "/create", nil)
	assert.Equal()
}

func TestRead(t *testing.T) {
	t.Skip("Not implemented")
}

func TestUpdate(t *testing.T) {
	t.Skip("Not implemented")
}

func TestDelete(t *testing.T) {
	t.Skip("Not implemented")
}
