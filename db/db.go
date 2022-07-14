package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var ctx = context.Background()

var (
	ARANGO_ROOT_PASSWORD = os.Getenv("ARANGO_ROOT_PASSWORD")
	DB_PROTOCOL          = os.Getenv("DB_PROTOCOL")
	DB_HOST              = os.Getenv("DB_HOST")
	DB_PORT              = os.Getenv("DB_PORT")
	DB_NAME              = os.Getenv("DB_NAME")
	DB_COLLECTION        = os.Getenv("DB_COLLECTION")
)

type DB struct {
	client   driver.Client
	conn     driver.Connection
	database driver.Database
	col      driver.Collection
	meta     driver.DocumentMeta // metadata of the last document created
	mutex    sync.Mutex
}

func (db *DB) Connect() {
	dbUrl := fmt.Sprintf("%s://%s:%s", DB_PROTOCOL, DB_HOST, DB_PORT)
	fmt.Printf("URL: %s\n", dbUrl)

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{dbUrl},
	})
	if err != nil {
		log.Fatal(err)
	}
	db.conn = conn
}

func (db *DB) Authenticate() {
	auth := driver.BasicAuthentication("root", ARANGO_ROOT_PASSWORD)
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     db.conn,
		Authentication: auth,
	})
	if err != nil {
		log.Fatal(err)
	}
	db.client = client
}

// creates a database if non-existant, initializes the driver.Database
func (db *DB) InitializeDatabase() {
	exists, err := db.client.DatabaseExists(ctx, DB_NAME)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err = db.client.CreateDatabase(ctx, DB_NAME, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	database, err := db.client.Database(ctx, DB_NAME)
	if err != nil {
		log.Fatal(err)
	}
	db.database = database
}

// creates a collection if non-existant, initializes the driver.Collection
func (db *DB) InitializeCollection() {
	exists, err := db.database.CollectionExists(ctx, DB_COLLECTION)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err = db.database.CreateCollection(ctx, DB_COLLECTION, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	col, err := db.database.Collection(ctx, "col")
	if err != nil {
		log.Fatal(err)
	}
	db.col = col
}

func (db *DB) UpdateMeta(meta driver.DocumentMeta) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.meta = meta
}

func (db *DB) Create(user User) {
	// check if record exists
	query := fmt.Sprintf(
		`FOR u IN col FILTER u.name == '%s' RETURN u`,
		user.Name,
	)
	cursor, err := db.col.Database().Query(
		ctx,
		query,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close()

	exists := cursor.Count() != 0

	// create record if not exists, read record
	if !exists {
		// create record
		meta, err := db.col.CreateDocument(ctx, user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("create res: %+v\n", meta)
		db.UpdateMeta(meta)
	}
}

// read record
func (db *DB) ReadAll() {
	cursor, err := db.col.Database().Query(
		ctx,
		`FOR u IN col FILTER u.name == 'Piotr' RETURN u`,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	for cursor.HasMore() {
		var user User
		meta, err := cursor.ReadDocument(ctx, &user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("read res: %+v\n", user)
		db.UpdateMeta(meta)
	}
}

func (db *DB) ReadOne() {
	cursor, err := db.col.Database().Query(
		ctx,
		`FOR u IN col FILTER u.name == 'Piotr' RETURN u`,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	var user User
	meta, err := cursor.ReadDocument(ctx, &user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read res: %+v\n", user)
	db.UpdateMeta(meta)
}

// update record
func (db *DB) Update() {
	patch := map[string]int{"Age": 21}
	meta, err := db.col.UpdateDocument(ctx, db.meta.Key, patch)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("update res: %+v\n", meta)
	db.UpdateMeta(meta)
}

// delete record
func (db *DB) Delete() {
	meta, err := db.col.RemoveDocument(ctx, db.meta.Key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("delete res: %+v\n", meta)
	db.UpdateMeta(meta)
}

func (db *DB) DeleteDB() {
	err := db.database.Remove(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
