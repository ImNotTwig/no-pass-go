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

func OpenAccountFromFile(account_path string) (Account, error) {
	account_path = ConvertToHashedPath(account_path)
	var account Account

	if _, err := os.Stat(account_path); err != nil {
		return account, err
	}

	decr_acc_file, err := exec.Command("./sh/decr.sh", account_path).Output()
	if err != nil {
		return account, err
	}

	if err = json.Unmarshal(decr_acc_file, &account); err != nil {
		return account, err
	}

	return account, nil
}

func RemoveAccount(account_path string) error {
	RemoveFromTreeFile(account_path)
	account_path = ConvertToHashedPath(account_path)

	if _, err := os.Stat(account_path); !os.IsNotExist(err) && err != nil {
		return err
	}

	os.Remove(account_path)

	return nil
}

func ConvertToHashedPath(account_path string) string {
	var hash []byte
	for _, block := range sha256.Sum256([]byte(account_path)) {
		hash = append(hash, block)
	}
	absolute_path, _ := filepath.Abs(config.BaseDirectory)
	return absolute_path + "/" + fmt.Sprintf("%x", hash)
}

// converts an account struct to json, encrypts it, and stores it to the given path
func SaveAccountToFile(account Account, account_path string) {
	AddToTreeFile(account_path)
	account_path = ConvertToHashedPath(account_path)

	// checking if the password file exists, and if it doesnt, we create it
	if _, err := os.Stat(account_path); os.IsNotExist(err) {
		os.Create(account_path)
	} else if err != nil {
		panic(err.Error())
	}

	json_data, err := json.Marshal(account)
	if err != nil {
		panic(err.Error())
	}
	// open and write the encrypted json data to the password file
	out, err := exec.Command("./sh/encr_string.sh", config.GPGPublicKey, string(json_data), account_path).Output()
	if err != nil {
		panic(err.Error())
	}
	os.WriteFile(account_path, out, os.ModePerm)
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

	gpg_encr_command := exec.Command(`./sh/encr_string.sh`, config.GPGPublicKey, treedb_string, absolute_path)
	gpg_encr_output, err := gpg_encr_command.Output()

	os.WriteFile(absolute_path, gpg_encr_output, os.ModePerm)

	return nil
}

func RemoveFromTreeFile(path string) error {
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

	delete(treedb, fmt.Sprintf("%x", path_hash))

	var treedb_string string
	for hash, path := range treedb {
		treedb_string += path + ":" + hash + "\n"
	}

	gpg_encr_command := exec.Command(`./sh/encr_string.sh`, config.GPGPublicKey, treedb_string, absolute_path)
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

	gpg_command := exec.Command(`./sh/decr.sh`, TreeDataBasePath)
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
