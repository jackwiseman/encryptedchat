package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

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
