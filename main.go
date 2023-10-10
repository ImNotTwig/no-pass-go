package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-json"
	// "golang.org/x/crypto/bcrypt"
)

const PASSWORD = "7y7*YA&*Y78y34y*&AYSy8ufuyhbdf^&teuyrgG&^DTFYUGWEHR87atyeruGRyuiwegr"

const baseDir = "/workspaces/no-pass-go/passwords"

const passwordData = `password: epicPassword
username: epicUsername
email: epicEmail
recovery_codes: [epicRecoveryCode0, epicRecoveryCode1]
domain: epicWebsite.com`

type Account struct {
	Password      string       `json:"Password"`
	Username      string       `json:"Username"`
	Email         string       `json:"Email"`
	Service       string       `json:"Service"`
	RecoveryCodes []string     `json:"RecoveryCodes"`
	ExtraData     *interface{} `json:"ExtraData"`
}

func main() {
	// hash the password
	// hash, err := bcrypt.GenerateFromPassword([]byte(PASSWORD), bcrypt.DefaultCost)
	// if err != nil {
	// 	panic(err.Error())
	// }

	account := Account{
		Password:      "34324324",
		Username:      "asdasdad",
		Email:         "something@something.com",
		Service:       "something.com",
		RecoveryCodes: []string{"asasd", "asdads"},
	}

	save_to_file(account, "passwords/email/google.com/Derpking37")

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

func save_to_file(account Account, password_path string) error {
	path_list := MakePathList(password_path)
	dirHash, fileHash := HashPathListIntoPath(path_list)

	if _, err := os.Stat(dirHash + fileHash); os.IsNotExist(err) {
		os.MkdirAll("passwords/"+dirHash, os.ModePerm)
		absPath, err := filepath.Abs("passwords/" + dirHash + fileHash)
		if err != nil {
			panic(err.Error())
		}
		_, err = os.Create(absPath)
		fmt.Println(dirHash + fileHash)
		if err != nil {
			panic(err.Error())
		}
	}

	file, _ := os.Create("passwords/" + dirHash + fileHash)

	jsonData, err := json.Marshal(passwordData)
	if err != nil {
		panic(err.Error())
	}

	encryptedJson, err := Encrypt([]byte(PASSWORD), jsonData)
	if err != nil {
		panic(err.Error())
	}

	file.Write(encryptedJson)

	return nil
}

func HashPathListIntoPath(path_list []string) (string, string) {
	var pathHash string
	for i := 0; i < len(path_list); i++ {
		var hash []byte
		for _, block := range sha256.Sum256([]byte(path_list[i])) {
			hash = append(hash, block)
		}
		pathHash += fmt.Sprintf("%x", hash) + "/"
	}
	var fileHash []byte
	for _, block := range sha256.Sum256([]byte(path_list[len(path_list)-1])) {
		fileHash = append(fileHash, block)
	}
	pathHash = strings.TrimSuffix(pathHash, fmt.Sprintf("%x", fileHash))
	return pathHash, fmt.Sprintf("%x", fileHash)
}

func MakePathList(password_file_path string) []string {
	var path_list []string
	var parent, absPath string

	for absPath != baseDir && absPath+"/" != baseDir {
		password_file_path = filepath.Dir(password_file_path)
		absPath, _ = filepath.Abs(password_file_path)
		parent = filepath.Base(password_file_path)
		path_list = append(path_list, parent)
	}
	for i := 0; i < len(path_list)/2; i++ {
		j := len(path_list) - i - 1
		path_list[i], path_list[j] = path_list[j], path_list[i]
	}
	return path_list[1:]
}
