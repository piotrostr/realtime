package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/piotrostr/realtime/rest"
)

func main() {
	router := rest.GetRouter()

	router.Run(":8080")
}
