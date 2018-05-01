package store

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

var DBPath = "yasssd.db"

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (a *Account) AllFiles() ([]string, error) {
	// find files/{username} bucket, iterate over contents
	return []string{"yep"}, nil
}

func PrintDump() {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("<<<db")
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("accounts"))
		if b != nil {
			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("accounts/%s: %s\n", k, v)
			}
		}

		b = tx.Bucket([]byte("tokens"))
		if b != nil {
			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("tokens/%s: %s\n", k, v)
			}
		}

		return nil
	})
	fmt.Println(">>>")

	return
}

func Truncate() {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("tokens"))
		tx.DeleteBucket([]byte("accounts"))
		return nil
	})

	return
}
