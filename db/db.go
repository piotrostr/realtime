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

// required env vars for establishing the conn
var (
	ARANGO_ROOT_PASSWORD = os.Getenv("ARANGO_ROOT_PASSWORD")
	DB_PROTOCOL          = os.Getenv("DB_PROTOCOL")
	DB_HOST              = os.Getenv("DB_HOST")
	DB_PORT              = os.Getenv("DB_PORT")
	DB_NAME              = os.Getenv("DB_NAME")
	DB_COLLECTION        = os.Getenv("DB_COLLECTION")
)

// main struct for the database
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
	// TODO return all of the errors and handle in the router to prevent
	// 50X errors
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

// update the meta on the database struct
func (db *DB) UpdateMeta(meta driver.DocumentMeta) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.meta = meta
}

// create a record
func (db *DB) Create(user User) *driver.DocumentMeta {
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
		return &meta
	}
	return nil
}

// read all records
func (db *DB) ReadAll() []*User {
	cursor, err := db.col.Database().Query(
		ctx,
		`FOR u IN col FILTER u.name == 'Piotr' RETURN u`,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	var users []*User
	for cursor.HasMore() {
		var user User
		meta, err := cursor.ReadDocument(ctx, &user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("read res: %+v\n", user)
		users = append(users, &user)
		db.UpdateMeta(meta)
	}
	return users
}

// read one record by name
func (db *DB) ReadOne(name string) (*User, *driver.DocumentMeta, error) {
	q := fmt.Sprintf(`FOR u IN col FILTER u.name == '%s' RETURN u`, name)
	cursor, err := db.col.Database().Query(ctx, q, nil)
	if err != nil {
		return nil, nil, err
	}
	var user User
	meta, err := cursor.ReadDocument(ctx, &user)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("read res: %+v\n", user)
	db.UpdateMeta(meta)
	return &user, &meta, err
}

// update record
func (db *DB) Update(user User) *driver.DocumentMeta {
	// check if record exists
	_, meta, err := db.ReadOne(user.Name)
	if err != nil {
		return nil
	}
	*meta, err = db.col.UpdateDocument(ctx, meta.Key, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("update res: %+v\n", meta)
	db.UpdateMeta(*meta)
	return meta
}

// delete record
func (db *DB) Delete(name string) *driver.DocumentMeta {
	// check if record exists
	_, meta, err := db.ReadOne(name)
	if err != nil {
		return nil
	}
	*meta, err = db.col.RemoveDocument(ctx, meta.Key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("delete res: %+v\n", meta)
	db.UpdateMeta(*meta)
	return meta
}

// delete the whole db
func (db *DB) DeleteDB() {
	err := db.database.Remove(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
