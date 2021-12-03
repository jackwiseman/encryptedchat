package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
	"os"
)

func genKeys() {
	reader := rand.Reader
	bits := 512
	key, err := rsa.GenerateKey(reader, bits)
	if err != nil {
		panic(err)
	}

	savePrivateKey("private.key", key)
}

func savePrivateKey(fileName string, key interface{}) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(key)
	if err != nil {
		panic(err)
	}
	file.Close()
}

// Usage
// var key rsa.PrivateKey
// loadPrivateKey(&key)
func loadPrivateKey(key interface{}) {
	// err needs to be handled a little differently here
	file, err := os.Open("private.key")
	if err != nil {
		panic(err)
	}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(key)
	if err != nil {
		panic(err)
	}
	file.Close()
}

