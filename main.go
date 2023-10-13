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

	_, err = toml.Decode(string(conf_file), &config)
	if err != nil {
		log.Fatalln("Invalid Config: ", err)
	}
	SaveAccountToFile(account, "email/tutanota.com/TwigTheCat")
}
