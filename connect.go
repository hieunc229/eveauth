package eveauth

import (
	"github.com/boltdb/bolt"
)

func getDB() (*bolt.DB, error) {
	return bolt.Open(AUTH_PATH, 0666, nil)
}
