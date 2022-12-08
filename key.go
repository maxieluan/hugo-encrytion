package main

import (
	"os"
	"crypto/sha256"
	"crypto/md5"
	"log"
	"time"
	"encoding/hex"
)

func return_key() string {
	key := os.Getenv("HUGO_ENCRYPTION_KEY")
	// check key, if not exist, raise error
	if key == "" {
		log.Fatal("HUGO_ENCRYPTION_KEY not set")
	}

	// append key with the UTC date
	key = key + time.Now().UTC().Format("2006-01-02")

	// sha256 of key
	hash := sha256.New()
	hash.Write([]byte(key))
	key = hex.EncodeToString(hash.Sum(nil))

	// md5 of the key
	hash = md5.New()
	hash.Write([]byte(key))
	key = hex.EncodeToString(hash.Sum(nil))

	return key
}