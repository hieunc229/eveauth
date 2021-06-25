package eveauth

import (
	"github.com/boltdb/bolt"
)

func getDB() (*bolt.DB, error) {
	return bolt.Open("store/auth", 0666, nil)
}
