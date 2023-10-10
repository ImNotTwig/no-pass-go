package main

import (
    "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
)

// Decrypts something
func Decrypt(key, data []byte) ([]byte, error) {
	var hash []byte
	for _, block := range sha256.Sum256(key) {
		hash = append(hash, block)
	}
	blockCipher, err := aes.NewCipher(hash)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Encrypts something
func Encrypt(key, data []byte) ([]byte, error) {
	var hash []byte
	for _, block := range sha256.Sum256(key) {
		hash = append(hash, block)
	}
	blockCipher, err := aes.NewCipher(hash)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}
