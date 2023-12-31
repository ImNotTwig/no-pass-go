# no-pass-go (WIP)
no-pass-go (npg) is a password/account manager written in Golang, very much inspired by [pass](https://www.passwordstore.org/), the name however is a reference to the Monopoly phrase "Do not pass Go and do not collect $200". 

With npg you have three _essential_ commands:
- show
- add
- remove

Other commmands are provided, but the above commands are all you need to use this program.

Filepath arguments given to this program should be relative to the account store folder if the filepath to the account file is `passwords/email/google.com/username`, then you only need to provide `email/google.com/username` as an argument.

- `npg show` shows the password, and optionally the metadata of a given account filepath. e.g: `email/google.com/username`
- `npg add` adds an account to the database, it requires a password to be defined. All the fields are defined through flags, including the password. a filepath is required.
- `npg remove` removes an account from the database, you only need to provide the filepath.

There are 3 major features of this program
1. It stores metadata for you, including usernames, emails, and service names (websites). It not only stores it for you, its stored in json, (encrypted of course)
2. The accounts are stored individually in their own files, which are hashes of the filepath given, the filepaths are stored in an index file called `pass_tree.asc`, each filepath has its own line, and is formatted as `path/to/account_file:hash_of_filepath`. The file is of course encrypted, as to not leak metadata
3. You can easily integrate this program into outside programs or scripts, because of the way that data is printed with the `show` command, through normal plaintext messages, or json

You need the go compiler to run this program. You can find out how to install Go on your system by using your favorite search engine and searching "how to install go on <operating_system>"

Make sure you rename `example_config.toml` to `config.toml`, and put in your **public** GPG key, and the absolute path to where you want to store your account data.
After making renaming and putting in your gpg key and account store directory in `config.toml`, move it to `~/.config/npg/config.toml`
> The `BaseDirectory` field needs an absolute path e.g: `/home/twig/passwords` **NOT** `~/passwords`

Here is a [guide](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/security_guide/sect-security_guide-encryption-gpg-creating_gpg_keys_using_the_command_line) on how to make a GPG key

to build and run the program after installing go and setting up your config file:
```sh
git clone https://github.com/ImNotTwig/no-pass-go
cd no-pass-go
go build
./npg
```
