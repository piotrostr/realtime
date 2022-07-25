package main

import (
	_ "github.com/joho/godotenv/autoload"
	database "github.com/piotrostr/realtime/db"
)

func main() {
	db := database.DB{}
	db.Connect()
	db.Authenticate()
	db.InitializeDatabase()
	db.InitializeCollection()
	db.Create(database.User{Name: "Piotr", Age: 30})
	db.ReadOne("Piotr")
	db.Update(database.User{Name: "Piotr", Age: 22})
	db.ReadOne("Piotr")
	db.Delete("Piotr")

	// clean up, skip in prod
	db.DeleteDB()
}
