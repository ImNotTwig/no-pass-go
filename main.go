package main

import (
	"fmt"
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
						Aliases:  []string{"p"},
						Value:    "",
						Required: true,
						Usage:    "The password you want to assign to this account",
					},
					&cli.StringFlag{
						Name:    "username",
						Aliases: []string{"u"},
						Value:   "",
						Usage:   "The username you want to assign to this account",
					},
					&cli.StringFlag{
						Name:    "email",
						Aliases: []string{"e"},
						Value:   "",
						Usage:   "The email you want to assign to this account",
					},
					&cli.StringFlag{
						Name:    "service",
						Aliases: []string{"s"},
						Value:   "",
						Usage:   "The service/website you want to assign to this account",
					},
				},
				Action: func(ctx *cli.Context) error {
					var account Account

					account_path := ctx.Args().First()
					if account_path == "" {
						return fmt.Errorf("No filepath given to store account.")
					}

					account.Password = ctx.Value("password").(string)
					account.Username = ctx.Value("username").(string)
					account.Email = ctx.Value("email").(string)
					account.Service = ctx.Value("service").(string)

					SaveAccountToFile(account, account_path)
					fmt.Println("Saved ", account_path, " to the store.")

					return nil
				},
			},
			{
				Name:      "remove",
				Aliases:   []string{"rm"},
				Usage:     "Remove an account from the store",
				UsageText: "npg rm, remove account_path",
				Action: func(ctx *cli.Context) error {
					account_path := ctx.Args().First()
					if err := RemoveAccount(account_path); err != nil {
						return err
					}
					fmt.Println("Removed ", account_path, " from the store.")
					return nil
				},
			},
			{
				Name: "show",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "Shows all account data, instead of just the password",
						Value:   false,
					},
				},
				Usage:     "Show an account's data from the store",
				UsageText: "npg show account_path",
				Action: func(ctx *cli.Context) error {
					account_path := ctx.Args().First()
					account, err := OpenAccountFromFile(account_path)
					if err != nil {
						return err
					}
					fmt.Println(account.Password)
					if ctx.Value("all").(bool) == false {
						return nil
					}

					fmt.Println("username: ", account.Username)
					fmt.Println("email: ", account.Email)
					fmt.Println("service: ", account.Service)

					return nil
				},
			},
			{
				Name:      "move",
				Aliases:   []string{"mv"},
				Usage:     "Move an account to another path",
				UsageText: "npg move OLD_PATH NEW_PATH",
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() < 2 {
						return fmt.Errorf("Not enough arguments. OLD_PATH and NEW_PATH are required.")
					}
					old_account_path := ctx.Args().First()
					new_account_path := ctx.Args().Get(1)
					account, err := OpenAccountFromFile(old_account_path)
					if err != nil {
						return err
					}

					RemoveAccount(old_account_path)
					SaveAccountToFile(account, new_account_path)

					fmt.Println("Moved ", old_account_path, " to ", new_account_path)

					return nil
				},
			},
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "List all the accounts in the store (without showing data associated).",
				UsageText: "npg list",
				Action: func(ctx *cli.Context) error {
					tree_database_string, err := ListTreeFile()
					if err != nil {
						return err
					}

					fmt.Println(tree_database_string)

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
}
