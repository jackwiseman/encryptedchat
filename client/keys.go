package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"os"
)

func genKeys() {
	reader := rand.Reader
	bits := 2048
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
		if _, ok := err.(*os.PathError); ok {
			genKeys()
			loadPrivateKey(key)
			return
		} else {
		panic(err)
		}
	}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(key)
	if err != nil {
		panic(err)
	}

	file.Close()
}

func encrypt(message string, key rsa.PublicKey) string {
	r := rand.Reader
	encrypted, err := rsa.EncryptOAEP(sha256.New(), r, &key, []byte(message), []byte("OAEP"))
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(encrypted)
}

func decrypt(message string, key rsa.PrivateKey) string {
	msgBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		panic(err)
	}
	r := rand.Reader
	decrypted, err2 := rsa.DecryptOAEP(sha256.New(), r, &key, msgBytes, []byte("OAEP"))
	if err2 != nil {
		panic(err2)
	}
	return string(decrypted)
}
