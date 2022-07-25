package db

import (
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"
)

func Test_CheckIfExistsWhenExists(t *testing.T) {
	db := DB{}
	db.Init()

	user := User{Name: "Piotr", Age: 30}
	db.Create(user)

	exists, meta := db.CheckIfExists(user)
	assert.True(t, exists)
	assert.NotNil(t, meta)

	err := db.DeleteDB()
	if err != nil {
		t.Error(err)
	}
}

func Test_CheckIfExistsWhenDoesntExist(t *testing.T) {
	db := DB{}
	db.Init()

	user := User{Name: "Piotr", Age: 30}
	exists, meta := db.CheckIfExists(user)

	assert.False(t, exists)
	assert.Nil(t, meta)
	err := db.DeleteDB()
	if err != nil {
		t.Error(err)
	}
}

func Test_E2E(t *testing.T) {
	db := DB{}
	db.Init()

	db.Create(User{Name: "Piotr", Age: 30})
	db.ReadOne("Piotr")
	db.Update(User{Name: "Piotr", Age: 22})
	user, meta, err := db.ReadOne("Piotr")
	assert.NoError(t, err)
	assert.Equal(t, user.Age, 22)
	assert.Equal(t, user.Name, "Piotr")
	assert.NotNil(t, meta)

	db.Delete(User{Name: "Piotr", Age: 22})

	err = db.DeleteDB()
	if err != nil {
		t.Error(err)
	}
}
