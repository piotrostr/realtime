package main

import (
	"context"
	"fmt"
	"log"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var ctx = context.Background()

var PASSWORD = "root"

func main() {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	if err != nil {
		log.Fatal(err)
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", "root"),
	})
	if err != nil {
		log.Fatal(err)
	}
	exists, err := client.DatabaseExists(ctx, "realtime")
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err = client.CreateDatabase(ctx, "realtime", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	db, err := client.Database(ctx, "realtime")
	if err != nil {
		log.Fatal(err)
	}

	collection, err := db.Collection(ctx, "collection")
	if err != nil {
		log.Fatal(err)
	}

	doc := map[string]interface{}{"name": "Piotr", "age": 22}
	meta, err := collection.CreateDocument(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(meta)
}
