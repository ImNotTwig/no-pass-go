package main

import (
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

const PASSWORD = "7y7*YA&*Y78y34y*&AYSy8ufuyhbdf^&teuyrgG&^"

var baseDir = "passwords"

func main() {
	// hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(PASSWORD), bcrypt.DefaultCost)
	if err != nil {
		panic(err.Error())
	}

	account := Account{
		Password:      "34324324",
		Username:      "asdasdad",
		Email:         "something@something.com",
		Service:       "something.com",
		RecoveryCodes: []string{"asasd", "asdads"},
	}

	baseDir, err = filepath.Abs(baseDir)
	if err != nil {
		panic(err.Error())
	}
	save_to_file([]byte(passwordHash), account, "passwords/email/tutanota.com/TwigTheCat")
	// absFile, err := filepath.Abs("./passwords/3049a1f8327e0215ea924b9e4e04cd4b0ff1800c74a536d9b81d3d8ced9994d3/82244417f956ac7c599f191593f7e441a4fafa20a4158fd52e154f1dc4c8ed92/d4c9d9027326271a89ce51fcaf328ed673f17be33469ff979e8ab8dd501e664f")
	// if err != nil {
	// panic(err.Error())
	// }

	// file, err := os.ReadFile(absFile)
	// if err != nil {
	// panic(err.Error())
	// }

	// var data Account
	// decryptedFile, err := Decrypt([]byte(passwordHash), file)
	// if err != nil {
	// panic(err.Error())
	// }
	// err = json.Unmarshal(decryptedFile, &data)
	// if err != nil {
	// panic(err.Error())
	// }
	// fmt.Println(data.Username)
}
