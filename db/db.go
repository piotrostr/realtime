package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/piotrostr/realtime/logger"
	"go.uber.org/zap"
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
	logger   *zap.SugaredLogger
}

func (db *DB) Init() {
	db.logger = logger.Get()

	err := db.Connect()
	if err != nil {
		db.logger.Fatalf(err.Error())
	}

	err = db.Authenticate()
	if err != nil {
		db.logger.Fatalf(err.Error())
	}

	err = db.InitializeDatabase()
	if err != nil {
		db.logger.Fatalf(err.Error())
	}

	err = db.InitializeCollection()
	if err != nil {
		db.logger.Fatalf(err.Error())
	}
}

func (db *DB) Connect() error {
	dbUrl := fmt.Sprintf("%s://%s:%s", DB_PROTOCOL, DB_HOST, DB_PORT)
	db.logger.Infof("URL: %s\n", dbUrl)

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{dbUrl},
	})
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

func (db *DB) Authenticate() error {
	auth := driver.BasicAuthentication("root", ARANGO_ROOT_PASSWORD)

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     db.conn,
		Authentication: auth,
	})
	if err != nil {
		return err
	}

	db.client = client
	return nil
}

// creates a database if non-existant, initializes the driver.Database
func (db *DB) InitializeDatabase() error {
	exists, err := db.client.DatabaseExists(ctx, DB_NAME)
	if err != nil {
		return err
	}

	if !exists {
		_, err = db.client.CreateDatabase(ctx, DB_NAME, nil)
		if err != nil {
			return err
		}
	}

	database, err := db.client.Database(ctx, DB_NAME)
	if err != nil {
		return err
	}
	db.database = database
	return nil
}

// creates a collection if non-existant, initializes the driver.Collection
func (db *DB) InitializeCollection() error {
	exists, err := db.database.CollectionExists(ctx, DB_COLLECTION)
	if err != nil {
		return err
	}

	if !exists {
		_, err = db.database.CreateCollection(ctx, DB_COLLECTION, nil)
		if err != nil {
			return err
		}
	}

	col, err := db.database.Collection(ctx, DB_COLLECTION)
	if err != nil {
		db.logger.Fatalf(err.Error())
	}
	db.col = col
	return nil
}

// update the meta on the database struct
func (db *DB) UpdateMeta(meta driver.DocumentMeta) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.meta = meta
}

func (db *DB) CheckIfExists(user User) bool {
	// check if record exists
	query := fmt.Sprintf(
		`FOR u IN %s FILTER u.name == '%s' RETURN u`,
		db.col.Name(),
		user.Name,
	)
	cursor, err := db.col.Database().Query(
		ctx,
		query,
		nil,
	)
	if err != nil {
		db.logger.Fatalf(err.Error())
	}
	defer cursor.Close()
	exists := cursor.Count() != 0
	return exists
}

// create a record
func (db *DB) Create(user User) *driver.DocumentMeta {
	// check if exists
	exists := db.CheckIfExists(user)

	// create record if not exists, read record
	if !exists {
		// create record
		meta, err := db.col.CreateDocument(ctx, user)
		if err != nil {
			db.logger.Fatalf(err.Error())
		}
		db.logger.Infof("create res: %+v\n", meta)
		db.UpdateMeta(meta)
		return &meta
	}
	return nil
}

// read all records
func (db *DB) ReadAll() []*User {
	q := fmt.Sprintf(`FOR u IN %s RETURN u`, db.col.Name())
	cursor, err := db.col.Database().Query(ctx, q, nil)
	if err != nil {
		db.logger.Fatalf(err.Error())
	}

	var users []*User
	for cursor.HasMore() {
		var user User
		meta, err := cursor.ReadDocument(ctx, &user)
		if err != nil {
			db.logger.Fatalf(err.Error())
		}
		db.logger.Infof("read res: %+v\n", user)
		users = append(users, &user)
		db.UpdateMeta(meta)
	}

	return users
}

// read one record by name
func (db *DB) ReadOne(name string) (*User, *driver.DocumentMeta, error) {
	exists := db.CheckIfExists(User{Name: name})
	if !exists {
		return nil, nil, fmt.Errorf("user does not exist")
	}

	// get the user
	query := fmt.Sprintf(
		`FOR u IN %s FILTER u.name == '%s' RETURN u`,
		db.col.Name(),
		name,
	)
	cursor, err := db.col.Database().Query(
		ctx,
		query,
		nil,
	)
	var user User
	meta, err := cursor.ReadDocument(ctx, &user)
	if err != nil {
		return nil, nil, err
	}
	db.logger.Infof("read res: %+v\n", user)
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
		db.logger.Fatalf(err.Error())
	}
	db.logger.Infof("update res: %+v\n", meta)
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
		db.logger.Fatalf(err.Error())
	}
	db.logger.Infof("delete res: %+v\n", meta)
	db.UpdateMeta(*meta)
	return meta
}

// delete the whole db
func (db *DB) DeleteDB() error {
	err := db.database.Remove(ctx)
	if err != nil {
		return err
	}
	return nil
}
