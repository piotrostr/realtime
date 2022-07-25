package db

import (
	"fmt"

	driver "github.com/arangodb/go-driver"
)

func (db *DB) CheckIfExists(user User) (bool, *driver.DocumentMeta) {
	// check if record exists
	query := fmt.Sprintf(
		`FOR u IN %s FILTER u.name == "%s" RETURN u`,
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
	exists := cursor.Count() > 0

	if exists {
		meta, err := cursor.ReadDocument(ctx, new(interface{}))
		if err != nil {
			db.logger.Fatalf(err.Error())
		}
		return exists, &meta
	}
	return exists, nil
}

// create a record
func (db *DB) Create(user User) *driver.DocumentMeta {
	// check if exists
	exists, meta := db.CheckIfExists(user)

	// create record if not exists, read record
	if !exists {
		// create record
		meta, err := db.col.CreateDocument(ctx, user)
		if err != nil {
			db.logger.Fatalf(err.Error())
		}
		db.logger.Infof("create res: %+v", meta)
		db.UpdateMeta(meta)
		return &meta
	}
	return meta
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
		users = append(users, &user)
		db.UpdateMeta(meta)
	}
	db.logger.Infof("read %d entries", len(users))

	return users
}

// read one record by name
func (db *DB) ReadOne(name string) (*User, *driver.DocumentMeta, error) {
	exists, meta := db.CheckIfExists(User{Name: name})
	db.logger.Infof("ReadOne exists: %t, meta: %+v\n", exists, meta)
	if !exists {
		return nil, nil, fmt.Errorf("user does not exist")
	}

	var user User
	readMeta, err := db.col.ReadDocument(ctx, meta.Key, &user)
	if err != nil {
		return nil, nil, err
	}
	db.logger.Infof("read one res: %+v", user)

	db.UpdateMeta(readMeta)
	return &user, meta, err
}

// update record
func (db *DB) Update(user User) *driver.DocumentMeta {
	exists, meta := db.CheckIfExists(user)
	if !exists {
		return nil
	}
	updatedMeta, err := db.col.UpdateDocument(ctx, meta.Key, user)
	if err != nil {
		db.logger.Fatalf(err.Error())
	}
	db.logger.Infof("update res: %+v", meta)
	db.UpdateMeta(updatedMeta)
	return &updatedMeta
}

// delete record
func (db *DB) Delete(user User) *driver.DocumentMeta {
	// check if record exists
	exists, meta := db.CheckIfExists(User{Name: "Piotr"})
	fmt.Printf("exists: %t, meta: %+v\n", exists, meta)
	if !exists {
		db.logger.Errorf(
			"user with name %s does not exist in %s %s",
			user.Name,
			db.database.Name(),
			db.col.Name(),
		)
		return nil
	}

	removedMeta, err := db.col.RemoveDocument(ctx, meta.Key)
	if err != nil {
		db.logger.Fatalf(err.Error())
	}
	db.logger.Infof("delete res: %+v", removedMeta)

	db.UpdateMeta(removedMeta)

	return &removedMeta
}

// delete the whole db
func (db *DB) DeleteDB() error {
	err := db.database.Remove(ctx)
	if err != nil {
		return err
	}
	return nil
}
