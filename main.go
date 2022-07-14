package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var ctx = context.Background()

var (
	ARANGO_ROOT_PASSWORD = os.Getenv("ARANGO_ROOT_PASSWORD")
	DB_PROTOCOL          = os.Getenv("DB_PROTOCOL")
	DB_HOST              = os.Getenv("DB_HOST")
	DB_PORT              = os.Getenv("DB_PORT")
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	dbUrl := fmt.Sprintf("%s://%s:%s", DB_PROTOCOL, DB_HOST, DB_PORT)
	fmt.Println(dbUrl)

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{dbUrl},
	})
	if err != nil {
		log.Fatal(err)
	}

	auth := driver.BasicAuthentication("root", ARANGO_ROOT_PASSWORD)
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: auth,
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

	exists, err = db.CollectionExists(ctx, "col")
	if err != nil {
		log.Fatal(err)
	}

	// check if collection exists, if not create it
	if !exists {
		_, err = db.CreateCollection(ctx, "col", nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	col, err := db.Collection(ctx, "col")
	if err != nil {
		log.Fatal(err)
	}

	user := User{
		Name: "Piotr",
		Age:  22,
	}

	// check if record exists
	cursor, err := col.Database().Query(
		ctx,
		`FOR u IN col FILTER u.name == 'Piotr' RETURN u`,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	exists = cursor.Count() != 0

	// create record if not exists, read record
	var meta driver.DocumentMeta
	if !exists {
		// create record
		meta, err = col.CreateDocument(ctx, user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("create res: %+v\n", meta)

		// read record
		var res User
		meta, err = col.ReadDocument(ctx, meta.Key, &res)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("read res: %+v\n", res)
	} else {
		// read record
		var res User
		meta, err = cursor.ReadDocument(ctx, &res)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("read res: %+v\n", res)
	}

	// update record
	user.Age = 23
	meta, err = col.UpdateDocument(ctx, meta.Key, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("update res: %+v\n", meta)

	// delete record
	meta, err = col.RemoveDocument(ctx, meta.Key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("delete res: %+v\n", meta)
}
