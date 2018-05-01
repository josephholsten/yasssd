package store

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/boltdb/bolt"
)

var alnumPattern = regexp.MustCompile("^[[:alnum:]]*$")

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *Account) ID() []byte {
	return []byte(a.Username)
}

func (a *Account) Marshal() ([]byte, error) {
	return json.Marshal(a)
}

func AccountUnmarshal(r io.Reader) (*Account, error) {
	var a Account
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&a)
	return &a, err
}

func (a *Account) IsValid() (ok bool, msg string) {
	if len(a.Username) < 3 {
		return false, "username too short, must be at least 3 characters"
	}
	if len(a.Username) > 20 {
		return false, "username too long, must be no more than 20 characters"
	}

	if !alnumPattern.MatchString(a.Username) {
		return false, "username contains invalid characters, must contain only alphanumeric characters"
	}

	if len(a.Password) < 8 {
		return false, "password too short, must be at least 8 characters"
	}

	return true, ""
}

func (a *Account) IsAuthenticated() bool {
	saved, err := FindAccountByID(a.ID())
	// If any error occurs, we don't trust the saved value and so reject authentication
	if err != nil {
		return false
	}
	// FIXME: this is unacceptably insecure. Ideally only a hashed password would be saved. If raw password must be saved, then perhaps the client could use CHAP handshake?
	if a.Username == saved.Username && a.Password == saved.Password {
		return true
	}
	return false
}

func (a *Account) Create() error {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("accounts"))
		if err != nil {
			return err
		}
		val, err := a.Marshal()
		if err != nil {
			return err
		}
		return b.Put(a.ID(), val)
	})
}

func FindAccountByID(id []byte) (*Account, error) {
	var a Account
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("accounts"))
		if b == nil {
			return fmt.Errorf("accounts bucket empty")
		}
		val := b.Get(id)
		if val == nil {
			return fmt.Errorf("account:%q not found", string(id))
		}
		return json.Unmarshal(val, &a)
	})
	if err != nil {
		return nil, err
	}
	return &a, nil
}

type TokenResponse struct {
	Token string `json:"token"`
}

func FindAccountByToken(t string) (*Account, error) {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var id []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tokens"))
		if b == nil {
			return fmt.Errorf("session token not recognized")
		}
		id = b.Get([]byte(t))
		if id == nil {
			return fmt.Errorf("session token not recognized")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	db.Close() // ensure closed before next db call
	return FindAccountByID(id)
}

func (a *Account) CreateToken() (*TokenResponse, error) {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var id string
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("tokens"))
		if err != nil {
			return err
		}

		rawID, _ := b.NextSequence()
		id = fmt.Sprintf("%d", rawID)

		// Persist bytes to users bucket.
		return b.Put([]byte(id), a.ID())
	})

	return &TokenResponse{id}, nil
}
