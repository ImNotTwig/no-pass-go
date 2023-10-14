package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

var config Config

func main() {
	conf_file, err := os.ReadFile("config.toml")
	if err != nil {
		panic(err.Error())
	}

	_, err = toml.Decode(string(conf_file), &config)
	if err != nil {
		log.Fatalln("Invalid Config: ", err)
	}

	app := &cli.App{
		Name: "npg",
		Commands: []*cli.Command{
			{
				Name:      "add",
				Usage:     "Add an account to the store",
				UsageText: "npg add [options] account_path",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "password",
						Value:    "",
						Required: true,
						Usage:    "The password you want to assign to this account",
					},
					&cli.StringFlag{
						Name:  "username",
						Value: "",
						Usage: "The username you want to assign to this account",
					},
					&cli.StringFlag{
						Name:  "email",
						Value: "",
						Usage: "The email you want to assign to this account",
					},
					&cli.StringFlag{
						Name:  "service",
						Value: "",
						Usage: "The service/website you want to assign to this account",
					},
				},
				Action: func(ctx *cli.Context) error {
					var account Account
					account_path := ctx.Args().First()
					account.Password = ctx.Value("password").(string)
					account.Username = ctx.Value("username").(string)
					account.Email = ctx.Value("email").(string)
					account.Service = ctx.Value("service").(string)

					SaveAccountToFile(account, account_path)

					return nil
				},
			},
		},
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Suggest:                true,
		Usage:                  "manage your password/account data",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	// account := Account{
	// 	Password:      "34324324",
	// 	Username:      "asdasdad",
	// 	Email:         "something@something.com",
	// 	Service:       "something.com",
	// 	RecoveryCodes: []string{"asasd", "asdads"},
	// }

	// SaveAccountToFile(account, "email/tutanota.com/TwigTheCat")
	// account, err = OpenAccountFromFile("email/tutanota.com/TwigTheCat")
	// if err != nil {
	// 	panic(err.Error())
	// }
}
