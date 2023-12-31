package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bitfield/script"
	"github.com/goccy/go-json"
)

// Read an account from file
func OpenAccountFromFile(account_path string) (Account, error) {
	account_path = ConvertToHashedPath(account_path)
	var account Account

	if _, err := os.Stat(account_path); err != nil {
		return account, err
	}

	decr_acc_file, err := script.Exec(fmt.Sprintf("gpg -dq %s", account_path)).String()

	if err = json.Unmarshal([]byte(decr_acc_file), &account); err != nil {
		return account, err
	}

	return account, nil
}

// edit account on file
func EditAccount(account_path string) error {
	editor := os.Getenv("EDITOR")
	if strings.TrimSpace(editor) == "" {
		return fmt.Errorf("Default editor not found, please set the EDITOR environmental variable.")
	}
	hashed_account_path := ConvertToHashedPath(account_path)

	temp_path, err := filepath.Abs(config.BaseDirectory + "/temp")
	if err != nil {
		return err
	}

	if _, err := os.Stat(temp_path); os.IsNotExist(err) {
		os.Create(temp_path)
	}

	old_account, err := script.Exec(fmt.Sprintf("gpg -dq %s", hashed_account_path)).String()
	if err != nil {
		return err
	}
	os.WriteFile(temp_path, []byte(old_account), os.ModePerm)

	cmd := exec.Command(editor, temp_path)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	decr_acc_file, err := os.ReadFile(temp_path)
	if err != nil {
		return err
	}

	var new_account Account
	if err = json.Unmarshal(decr_acc_file, &new_account); err != nil {
		return err
	}

	SaveAccountToFile(new_account, account_path)

	return nil
}

// Remove an account from the store
func RemoveAccount(account_path string) error {
	RemoveFromTreeFile(account_path)
	account_path = ConvertToHashedPath(account_path)

	if _, err := os.Stat(account_path); !os.IsNotExist(err) && err != nil {
		return err
	}

	os.Remove(account_path)

	return nil
}

// convert a filepath to a hash
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
	out, err := script.Echo(string(json_data)).Exec(fmt.Sprintf("gpg --encrypt --armor --always-trust --batch --yes --recipient %s --output -", config.GPGPublicKey)).String()
	if err != nil {
		panic(err.Error())
	}
	os.WriteFile(account_path, []byte(out), os.ModePerm)
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

	gpg_encr_output, err := script.Echo(string(treedb_string)).Exec(fmt.Sprintf("gpg --encrypt --armor --always-trust --batch --yes --recipient %s --output -", config.GPGPublicKey)).String()
	if err != nil {
		panic(err.Error())
	}

	os.WriteFile(absolute_path, []byte(gpg_encr_output), os.ModePerm)

	return nil
}

// remove a path from the pass_tree file
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

	gpg_encr_output, err := script.Echo(string(treedb_string)).Exec(fmt.Sprintf("gpg --encrypt --armor --always-trust --batch --yes --recipient %s --output -", config.GPGPublicKey)).String()

	os.WriteFile(absolute_path, []byte(gpg_encr_output), os.ModePerm)

	return nil
}

func ListTreeFile() (string, error) {
	tree_data_base, err := ParseTreeFile()
	if err != nil {
		return "", err
	}

	var treedb_string string
	for hash, path := range tree_data_base {
		treedb_string += path + ":" + hash + "\n"
	}

	return treedb_string, nil
}

// parse the tree database into a dictionary
func ParseTreeFile() (TreeDataBase, error) {
	TreeDataBasePath, _ := filepath.Abs(config.BaseDirectory + "/" + "pass_tree.asc")
	treedb := make(map[string]string)

	tree_file, _ := os.ReadFile(TreeDataBasePath)
	if string(tree_file) == "" {
		tree_file, _ = os.ReadFile(TreeDataBasePath)
	}

	if _, err := os.Stat(TreeDataBasePath); os.IsNotExist(err) {
		os.Create(TreeDataBasePath)
	}

	gpg_output, err := script.Exec(fmt.Sprintf("gpg -dq %s", TreeDataBasePath)).String()
	if err != nil {
		return treedb, nil
	}
	gpg_output = strings.TrimSpace(gpg_output)

	if gpg_output == "" {
		return treedb, nil
	}
	for _, data := range strings.Split(gpg_output, "\n") {
		PathAndHash := strings.Split(data, ":")
		treedb[PathAndHash[1]] = PathAndHash[0]
	}

	return treedb, nil
}
