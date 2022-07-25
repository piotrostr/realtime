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
	db.logger.Infof("URL: %s", dbUrl)

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
	db.logger.Infof("database: %+v", database)
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
