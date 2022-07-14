package main

import (
	"context"
	"fmt"
	"log"
	"os"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var ctx = context.Background()

var (
	ARANGO_ROOT_PASSWORD = os.Getenv("ARANGO_ROOT_PASSWORD")
	DB_HOST              = os.Getenv("DB_HOST")
	DB_PORT              = os.Getenv("DB_PORT")
)

func main() {
	dbUrl := fmt.Sprintf("http://%s:%s", DB_HOST, DB_PORT)
	fmt.Println(dbUrl)

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{dbUrl},
	})
	if err != nil {
		log.Fatal(err)
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", ARANGO_ROOT_PASSWORD),
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
