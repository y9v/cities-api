package test_helpers

import (
	"github.com/boltdb/bolt"
	"testing"
)

func CreateDB(t *testing.T) *bolt.DB {
	db, err := bolt.Open(CreateTempfile("", t), 0600, nil)
	if err != nil {
		t.Fatalf("DB: %v", err)
	}

	return db
}

func CreateBucket(t *testing.T, db *bolt.DB, name []byte) *bolt.Bucket {
	var b *bolt.Bucket

	err := db.Update(func(tx *bolt.Tx) error {
		var err error
		b, err = tx.CreateBucket(name)
		return err
	})
	if err != nil {
		t.Fatalf("DB: %v", err)
	}

	return b
}

func PutToBucket(
	t *testing.T, db *bolt.DB, bucketName []byte, key string, data string,
) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err := b.Put([]byte(key), []byte(data))
		return err
	})
	if err != nil {
		t.Fatalf("DB: %v", err)
	}
}

func ReadFromBucket(
	t *testing.T, db *bolt.DB, bucketName []byte, key string,
) string {
	var val []byte

	db.View(func(tx *bolt.Tx) error {
		if bucket := tx.Bucket(bucketName); bucket != nil {
			val = bucket.Get([]byte(key))
		}

		return nil
	})

	return string(val)
}

func DeleteFromBucket(t *testing.T, db *bolt.DB, bucketName []byte, key string) {
	db.View(func(tx *bolt.Tx) error {
		if bucket := tx.Bucket(bucketName); bucket != nil {
			bucket.Delete([]byte(key))
		}

		return nil
	})
}
