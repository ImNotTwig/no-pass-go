package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goccy/go-json"
)

func OpenAccountFromFile(key []byte, account_path string) {
}

// converts an account struct to json, encrypts it, and stores it to the given path
func SaveAccountToFile(account Account, account_path string) {
	var hash []byte
	for _, block := range sha256.Sum256([]byte(account_path)) {
		hash = append(hash, block)
	}
	absolute_path, _ := filepath.Abs("passwords/" + fmt.Sprintf("%x", hash))

	// checking if the password file exists, and if it doesnt, we create it
	if _, err := os.Stat(absolute_path); os.IsNotExist(err) {
		os.Create(absolute_path)
	} else if err != nil {
		panic(err.Error())
	}

	json_data, err := json.Marshal(account)
	if err != nil {
		panic(err.Error())
	}
	// open and write the encrypted json data to the password file
	out, err := exec.Command("./bash/encr_string.sh", config.GPGPublicKey, string(json_data), absolute_path).Output()
	if err != nil {
		panic(err.Error())
	}
	os.WriteFile(absolute_path, out, os.ModePerm)
	AddToTreeFile(account_path)
}

// add an entry to the tree database file
func AddToTreeFile(path string) error {
	absolute_path, _ := filepath.Abs(config.BaseDirectory + "/" + "pass_tree.asc")

	// checking if the database file exists
	if _, err := os.Stat(absolute_path); os.IsNotExist(err) {
		os.Create(absolute_path)
	}
	treedb, err := ParseTreeFile()
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
func ParseTreeFile() (TreeDataBase, error) {
	TreeDataBasePath, _ := filepath.Abs(config.BaseDirectory + "/" + "pass_tree.asc")
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
