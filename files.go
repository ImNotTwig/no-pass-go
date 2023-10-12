package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goccy/go-json"
)

// converts an account struct to json, encrypts it, and stores it to the given path
func SavePasswordToFile(account Account, password_path string) {
	absolute_path, _ := filepath.Abs("passwords/" + password_path)

	var hash []byte
	for _, block := range sha256.Sum256([]byte(password_path)) {
		hash = append(hash, block)
	}

	// checking if the password file exists, and if it doesnt, we create it
	if _, err := os.Stat(absolute_path); os.IsNotExist(err) {
		os.Create(absolute_path)

		// creating the password file
		_, err = os.Create(absolute_path)
	} else {
		log.Fatalln(err)
	}

	json_data, err := json.Marshal(account)
	if err != nil {
		log.Fatalln(err)
	}

	// open and write the encrypted json data to the password file
	password_file, _ := os.Create(absolute_path)
	password_file.Write(json_data)
	exec.Command("gpg", "--always-trust", "--batch", "--yes", "--encrypt", "--recipient", config.GPGPublicKey, absolute_path)
}

// add an entry to the tree database file
func AddToTreeFile(path string) error {
	absolute_path, _ := filepath.Abs(config.BaseDirectory + "/" + "pass_tree.asc")

	// checking if the database file exists
	if _, err := os.Stat(absolute_path); os.IsNotExist(err) {
		os.Create(absolute_path)
	}
	treedb, err := ParseTreeFile(absolute_path)
	if err != nil {
		return err
	}

	var path_hash []byte
	for _, block := range sha256.Sum256([]byte(path)) {
		path_hash = append(path_hash, block)
	}

	treedb[fmt.Sprintf("%x", path_hash)] = path

	var treedb_string string
	for hash, path := range treedb {
		treedb_string += path + ":" + hash + "\n"
	}

	gpg_encr_command := exec.Command(`./bash/encr_string.sh`, config.GPGPublicKey, treedb_string, absolute_path)
	gpg_encr_output, err := gpg_encr_command.Output()

	os.WriteFile(absolute_path, gpg_encr_output, os.ModePerm)

	return nil
}

// parse the tree database into a dictionary
func ParseTreeFile(TreeDataBasePath string) (TreeDataBase, error) {
	treedb := make(map[string]string)

	tree_file, _ := os.ReadFile(TreeDataBasePath)
	if string(tree_file) == "" {
		tree_file, _ = os.ReadFile(TreeDataBasePath)
	}

	gpg_command := exec.Command(`./bash/decr.sh`, TreeDataBasePath)
	gpg_output, err := gpg_command.Output()
	if err != nil {
		panic(err.Error())
	}
	gpg_output = []byte(strings.TrimSpace(string(gpg_output)))

	if string(gpg_output) == "" {
		return treedb, nil
	}
	for _, data := range strings.Split(string(gpg_output), "\n") {
		PathAndHash := strings.Split(data, ":")
		treedb[PathAndHash[1]] = PathAndHash[0]
	}

	return treedb, nil
}
