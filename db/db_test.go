package db

import (
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

func TestE2E(t *testing.T) {
	db := DB{}

	err := db.Connect()
	if err != nil {
		t.Error(err)
	}

	err = db.Authenticate()
	if err != nil {
		t.Error(err)
	}

	err = db.InitializeDatabase()
	if err != nil {
		t.Error(err)
	}

	err = db.InitializeCollection()
	if err != nil {
		t.Error(err)
	}

	db.Create(User{Name: "Piotr", Age: 30})
	db.ReadOne("Piotr")
	db.Update(User{Name: "Piotr", Age: 22})
	db.ReadOne("Piotr")
	db.Delete("Piotr")

	// clean up, skip in prod
	err = db.DeleteDB()
	if err != nil {
		t.Error(err)
	}
}
