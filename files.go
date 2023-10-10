package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-json"
)

// converts an account struct to json, encrypts it, and stores it to the given path
func save_to_file(key []byte, account Account, password_path string) error {
	path_list := MakePathList(password_path)
	dirHash, fileHash := HashPathListIntoString(path_list)

	if _, err := os.Stat(dirHash + fileHash); os.IsNotExist(err) {
		os.MkdirAll("passwords/"+dirHash, os.ModePerm)
		absPath, err := filepath.Abs("passwords/" + dirHash + fileHash)
		if err != nil {
			panic(err.Error())
		}
		_, err = os.Create(absPath)
		if err != nil {
			panic(err.Error())
		}
	}

	file, _ := os.Create("passwords/" + dirHash + fileHash)

	jsonData, err := json.Marshal(account)
	if err != nil {
		panic(err.Error())
	}

	encryptedJson, err := Encrypt(key, jsonData)
	if err != nil {
		panic(err.Error())
	}

	file.Write(encryptedJson)

	return nil
}

// takes a normal path, and hashes every directory, and returns the hashed path, as well as the hashed file at the end of the path
func HashPathListIntoString(path_list []string) (string, string) {
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
	pathHash = strings.TrimSuffix(pathHash, fmt.Sprintf("%x/", fileHash))
	return pathHash, fmt.Sprintf("%x", fileHash)
}

// returns a list of directories and the file from a path
func MakePathList(password_file_path string) []string {
	var path_list []string
	var parent, absPath string

	for absPath != baseDir && absPath+"/" != baseDir {
		parent = filepath.Base(password_file_path)
		password_file_path = filepath.Dir(password_file_path)
		absPath, _ = filepath.Abs(password_file_path)
		path_list = append(path_list, parent)
	}
	for i := 0; i < len(path_list)/2; i++ {
		j := len(path_list) - i - 1
		path_list[i], path_list[j] = path_list[j], path_list[i]
	}
	return path_list
}
