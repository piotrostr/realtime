package db

import (
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

func TestE2E(t *testing.T) {
	db := DB{}
	db.Connect()
	db.Authenticate()
	db.InitializeDatabase()
	db.InitializeCollection()
	db.Create(User{Name: "Piotr", Age: 30})
	db.ReadOne("Piotr")
	db.Update(User{Name: "Piotr", Age: 22})
	db.ReadOne("Piotr")
	db.Delete("Piotr")

	// clean up, skip in prod
	db.DeleteDB()
}
