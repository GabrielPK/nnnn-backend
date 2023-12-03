package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"golang.org/x/crypto/scrypt"
)

func HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Generate the hashed password
	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}

	// Return the salt and the hashed password, encoded in base64 and concatenated
	return base64.StdEncoding.EncodeToString(salt) + base64.StdEncoding.EncodeToString(dk), nil
}

func ComparePasswords(hashedPwd, plainPwd string) error {
	// Decode the salt (first 24 characters after base64 encoding of 16 bytes)
	salt, err := base64.StdEncoding.DecodeString(hashedPwd[:24])
	if err != nil {
		return err
	}

	// Decode the stored hash (the rest of the string)
	storedHash, err := base64.StdEncoding.DecodeString(hashedPwd[24:])
	if err != nil {
		return err
	}

	// Hash the provided password using the same salt
	hash, err := scrypt.Key([]byte(plainPwd), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return err
	}

	// Compare the hashes
	if !bytes.Equal(hash, storedHash) {
		return fmt.Errorf("password does not match")
	}
	return nil
}
