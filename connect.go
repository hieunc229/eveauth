package eveauth

import (
	"github.com/boltdb/bolt"
)

func GetDB() (*bolt.DB, error) {
	return bolt.Open("store/auth", 0666, nil)
}
