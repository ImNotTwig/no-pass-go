package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-json"
)

// converts an account struct to json, encrypts it, and stores it to the given path
func SaveToFile(key []byte, account Account, password_path string) error {
	path_list := MakePathList(password_path)
	dir_hash, file_hash := HashPathListIntoString(path_list)

	// checking if the password file path exists, and if it doesnt, we create it
	if _, err := os.Stat(dir_hash + file_hash); os.IsNotExist(err) {
		// making the path without the password file
		os.MkdirAll("passwords/"+dir_hash, os.ModePerm)

		// creating the password file
		absolute_path, _ := filepath.Abs("passwords/" + dir_hash + file_hash)
		_, err = os.Create(absolute_path)
	} else {
		log.Fatalln(err)
	}

	json_data, err := json.Marshal(account)
	if err != nil {
		log.Fatalln(err)
	}
	encrypted_json, err := Encrypt(key, json_data)
	if err != nil {
		log.Fatalln(err)
	}

	// open and write the encrypted json data to the password file
	password_file, _ := os.Create("passwords/" + dir_hash + file_hash)
	password_file.Write(encrypted_json)

	return nil
}

// takes a normal path, and hashes every directory, and returns the hashed path, as well as the hashed file at the end of the path
func HashPathListIntoString(path_list []string) (string, string) {
	var path_hash string
	var file_hash []byte
	file_name := path_list[len(path_list)-1]

	// for every file/directory in the pathlist, hash it, and convert the 32 byte slice, to a dynamic slice
	// skip the last element in the path_list slice, because that is the account file, and not a directory
	for i := 0; i < len(path_list)-1; i++ {
		var hash []byte
		dir_name := []byte(path_list[i])

		// converting the 32 byte slice to a dynamic slice
		for _, block := range sha256.Sum256(dir_name) {
			hash = append(hash, block)
		}

		// add the hashed directory name followed by a "/" to the path_hash
		path_hash += fmt.Sprintf("%x", hash) + "/"
	}

	// hashing the account file name
	for _, block := range sha256.Sum256([]byte(file_name)) {
		file_hash = append(file_hash, block)
	}

	return path_hash, fmt.Sprintf("%x", file_hash)
}

// returns a list of directories and the file from a path
func MakePathList(password_file_path string) []string {
	var path_list []string
	var last_element, absolute_path string

	for absolute_path != base_dir && absolute_path+"/" != base_dir {
		// Base returns the last file/directory of the path
		last_element = filepath.Base(password_file_path)

		// Dir returns the path without the last file/directory
		password_file_path = filepath.Dir(password_file_path)

		// We need the absolute path to test when we are in the base pasword directory or not
		absolute_path, _ = filepath.Abs(password_file_path)

		// Appending the last file/directory to the path_list
		path_list = append(path_list, last_element)
	}
	// reverse the list, because the path will be backwards
	for i := 0; i < len(path_list)/2; i++ {
		j := len(path_list) - i - 1
		path_list[i], path_list[j] = path_list[j], path_list[i]
	}
	return path_list
}

func AddToTreeFile(key []byte, path string) error {
	// we have to check if either the base_dir + / or if base_dir by itself is a valid path directory with / at the end
	absolute_path, err := filepath.Abs(base_dir + "TreeFile.db")
	if err != nil {
		absolute_path, err = filepath.Abs(base_dir + "/" + "TreeFile.db")
		if err != nil {
			return err
		}
	}
	// checking if the database file exists
	if _, err := os.Stat(absolute_path); os.IsNotExist(err) {
		os.Create(absolute_path)
	}
	treedb, err := ParseTreeFile(key, absolute_path)
	if err != nil {
		return err
	}

	var path_hash []byte
	for _, block := range sha256.Sum256([]byte(absolute_path)) {
		path_hash = append(path_hash, block)
	}

	treedb[fmt.Sprintf("%x", path_hash)] = absolute_path

	os.Create(absolute_path)

	treedb_file, _ := os.OpenFile(absolute_path, os.O_APPEND|os.O_WRONLY, 0666)
	defer treedb_file.Close()

	var treedb_string string

	for hash, path := range treedb {
		fmt.Println(hash)
		fmt.Println(path)

		treedb_string += path + ":" + hash + "\n"
	}

	encrypted_treedb, err := Encrypt(key, []byte(treedb_string))
	if err != nil {
		return err
	}

	treedb_file.Write([]byte(encrypted_treedb))

	return nil
}

func ParseTreeFile(key []byte, TreeDataBasePath string) (TreeDataBase, error) {
	treedb := make(map[string]string)

	tree_file, _ := os.ReadFile(TreeDataBasePath)
	if string(tree_file) == "" {
		fmt.Println("empty")
		return treedb, nil
	}
	var decrypted_database []byte
	if !(string(tree_file) == "") {
		var err error
		decrypted_database, err = Decrypt(key, tree_file)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("adding paths and hashes to dict")
	for _, data := range strings.Split(string(decrypted_database), "\n") {
		PathAndHash := strings.Split(data, ":")
		treedb[PathAndHash[1]] = PathAndHash[0]
	}

	return treedb, nil
}
