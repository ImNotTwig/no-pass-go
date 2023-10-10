package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
    "crypto/sha256"
    "golang.org/x/crypto/bcrypt"
	"fmt"
)

const PASSWORD = "7y7*YA&*Y78y34y*&AYSy8ufuyhbdf^&teuyrgG&^DTFYUGWEHR87atyeruGRyuiwegr"

func main() {
	// hash the password
    hash, err := bcrypt.GenerateFromPassword([]byte(PASSWORD), bcrypt.DefaultCost)
    fmt.Println(string(hash))

    passwordData := `password: epicPassword
username: epicUsername
email: epicEmail
recovery_codes: [epicRecoveryCode0, epicRecoveryCode1]
domain: epicWebsite.com`

    cipherText, err := Encrypt(hash, []byte(passwordData))
    if err != nil {
        panic(err.Error())
    }
    decryptedCipherText, err := Decrypt(hash, cipherText)
    if err != nil {
        panic(err.Error())
    }
    fmt.Println(string(decryptedCipherText))

}

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
