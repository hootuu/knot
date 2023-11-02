package dbx

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/hootuu/utils/errors"
	"github.com/hootuu/utils/logger"
	"go.uber.org/zap"
)

type Collection struct {
	db   *DBX
	name string
}

func (coll *Collection) Put(key string, data interface{}) *errors.Error {
	byteData, nErr := json.Marshal(data)
	if nErr != nil {
		return errors.Verify("json marshal data failed: " + nErr.Error())
	}
	nErr = coll.db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(coll.name))
		nErr := bucket.Put([]byte(key), byteData)
		if nErr != nil {
			return nErr
		}
		return nil
	})
	if nErr != nil {
		logger.Logger.Error("dbx.update error", zap.Error(nErr))
		return errors.Sys("dbx.update failed: " + nErr.Error())
	}
	return nil
}

func (coll *Collection) Get(key string) ([]byte, *errors.Error) {
	var byteData []byte = nil
	nErr := coll.db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(coll.name))
		byteData = bucket.Get([]byte(key))
		return nil
	})
	if nErr != nil {
		return nil, errors.Sys("dbx.load.data failed: " + nErr.Error())
	}
	return byteData, nil
}

func (coll *Collection) MustGet(key string, v interface{}) *errors.Error {
	byteData, err := coll.Get(key)
	if err != nil {
		return err
	}
	if byteData == nil {
		return errors.Verify("no such data: " + key)
	}
	nErr := json.Unmarshal(byteData, v)
	if nErr != nil {
		return errors.Sys("json.unmarshal data failed: " + nErr.Error())
	}
	return nil
}

func (coll *Collection) Iter(callback func(k string, byteData []byte) bool) *errors.Error {
	nErr := coll.db.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(coll.name))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			kStr := string(k)
			bNext := callback(kStr, v)
			if !bNext {
				return nil
			}
		}
		return nil
	})
	if nErr != nil {
		return errors.Sys("db.collection.iter field: " + nErr.Error())
	}
	return nil
}

func (coll *Collection) Remove(key string) *errors.Error {
	nErr := coll.db.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(coll.name))
		nErr := b.Delete([]byte(key))
		if nErr != nil {
			return nErr
		}
		return nil
	})
	if nErr != nil {
		return errors.Sys("dbx.delete failed: " + nErr.Error())
	}
	return nil
}

func doNewCollection(db *DBX, name string) (*Collection, *errors.Error) {
	if db == nil {
		return nil, errors.Verify("require dbx")
	}
	if !db.IsReady() {
		return nil, errors.Verify("dbx not ready")
	}
	if err := NameVerify(name); err != nil {
		return nil, err
	}
	tx, nErr := db.boltDB.Begin(true)
	if nErr != nil {
		return nil, errors.Sys("dbx.tx.Begin failed: " + nErr.Error())
	}
	defer func() {
		nErr := tx.Rollback()
		if nErr != nil {
			logger.Logger.Error("tx.Rollback() failed", zap.Error(nErr))
		}
	}()

	_, nErr = tx.CreateBucketIfNotExists([]byte(db.name))
	if nErr != nil {
		logger.Logger.Error("tx.CreateBucketIfNotExists([]byte(db.name)) error", zap.Error(nErr))
		return nil, errors.Sys("tx.CreateBucketIfNotExists failed: " + nErr.Error())
	}

	if nErr := tx.Commit(); nErr != nil {
		logger.Logger.Error("tx.Commit() error", zap.Error(nErr))
		return nil, errors.Sys("tx.Commit Failed: " + nErr.Error())
	}

	return &Collection{
		db:   db,
		name: name,
	}, nil
}
