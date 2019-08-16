/*
	Bolt DB Handling functions
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

var world = []byte("world")

func main() {
	db, err := BoltOpen("sample.db", nil)
	if err != nil {
		log.Fatalf("%v %T\n", err, err)
	}
	defer BoltClose(db)

	key := []byte("hello")
	value := []byte("Hello World!")

	err = BoltPutItem(db, world, key, value)
	value, err = BoltGetItem(db, world, key)

	key = []byte("stoney")
	value = []byte("Hello Kang!")

	err = BoltPutItem(db, world, key, value)
	value, err = BoltGetItem(db, world, key)

	key = []byte("hello")
	value = []byte("Hello Kang!")

	err = BoltPutItem(db, world, key, value)
	value, err = BoltGetItem(db, world, key)

	//err = BoltDeleteItem(db, world, key)
	//value, err = BoltGetItem(db, world, key)

	world = []byte("dark world")
	err = BoltPutItem(db, world, key, value)

	key = []byte("hello2")
	value = []byte("Hello Kang2!")
	err = BoltPutItem(db, world, key, value)

	err = BoltListBucket(db, world)
	err = BoltDeleteBucket(db, world)
	err = BoltListAll(db)

	BoltState(db)
	BoltMonitor(db, 5*time.Second)
}

func BoltOpen(dbpath string, dbopt *bolt.Options) (*bolt.DB, error) {
	return bolt.Open(dbpath, 0644, dbopt)
}

func BoltClose(db *bolt.DB) {
	db.Close()
}

func BoltMonitor(db *bolt.DB, ts time.Duration) {
	// Grab the initial stats.
	prev := db.Stats()

	for {
		// Wait for 10s.
		time.Sleep(ts)

		// Grab the current stats and diff them.
		stats := db.Stats()
		diff := stats.Sub(&prev)

		// Encode stats to JSON and print to STDERR.
		json.NewEncoder(os.Stderr).Encode(diff)

		// Save stats for the next loop.
		prev = stats
	}
}

func BoltState(db *bolt.DB) {
	stat := db.Stats()
	json.NewEncoder(os.Stderr).Encode(stat)
}

func BoltListAll(db *bolt.DB) error {
	err := db.View(func(tx *bolt.Tx) error {
		err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Printf("bk:%q\n", name)
			BoltListBucket(db, name)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatalf("%v %T\n", err, err)
	}

	return nil
}

func BoltListBucket(db *bolt.DB, bucket []byte) error {
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucket).Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			fmt.Printf("\tk:%q, v:%q\n", k, v)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("%v %T\n", err, err)
	}

	return nil
}

func BoltDeleteBucket(db *bolt.DB, bucket []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucket)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatalf("%v %T\n", err, err)
	}

	return nil
}

func BoltGetItem(db *bolt.DB, bucket, key []byte) ([]byte, error) {
	var value []byte

	err := db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return fmt.Errorf("Bucket %q not found!", bk)
		}

		value = bk.Get(key)
		if value == nil {
			log.Printf("Key %q not found\n", key)
		}
		fmt.Println(string(value))

		return nil
	})

	if err != nil {
		log.Fatalf("%v %T\n", err, err)
	}

	return value, err
}

func BoltPutItem(db *bolt.DB, bucket, key, value []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		err = bk.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatalf("%v %T\n", err, err)
	}

	return nil
}

func BoltDeleteItem(db *bolt.DB, bucket, key []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket(bucket)
		if bk == nil {
			return fmt.Errorf("Bucket %q not found!", bk)
		}

		err := bk.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatalf("%v %T\n", err, err)
	}

	return nil
}
