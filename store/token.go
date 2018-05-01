package store

import "github.com/boltdb/bolt"

func CreateToken(tokenID, accountID string) error {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("tokens"))
		if err != nil {
			return err
		}

		return b.Put([]byte(tokenID), []byte(accountID))
	})
}
