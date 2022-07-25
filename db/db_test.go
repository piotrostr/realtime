package db

import (
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

func TestE2E(t *testing.T) {
	db := DB{}
	db.Init()

	db.Create(User{Name: "Piotr", Age: 30})
	db.ReadOne("Piotr")
	db.Update(User{Name: "Piotr", Age: 22})
	db.ReadOne("Piotr")
	db.Delete(User{Name: "Piotr"})

	// clean up, skip in prod
	// err := db.DeleteDB()
	// if err != nil {
	// 	t.Error(err)
	// }
}
