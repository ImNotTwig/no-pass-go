package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var config Config

func main() {
	account := Account{
		Password:      "34324324",
		Username:      "asdasdad",
		Email:         "something@something.com",
		Service:       "something.com",
		RecoveryCodes: []string{"asasd", "asdads"},
	}

	conf_file, err := os.ReadFile("config.toml")
	if err != nil {
		panic(err.Error())
	}

	_, err = toml.Decode(string(conf_file), config)
	config.BaseDirectory = "passwords"
	config.GPGPublicKey = "6C75659CD4F90EB0E718C55EC98187AA4A03B7A0"

	SavePasswordToFile(account, "passwords/email/tutanota.com/TwigTheCat")
	err = AddToTreeFile("email/tutanota.com/ghostmage42")
	err = AddToTreeFile("email/tutanota.com/TwigTheCat")
	if err != nil {
		log.Fatalln(err)
	}
}
