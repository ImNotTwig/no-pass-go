package main

type Account struct {
	Password  string       `json:"Password"`
	Username  string       `json:"Username"`
	Email     string       `json:"Email"`
	Service   string       `json:"Service"`
	ExtraData *interface{} `json:"ExtraData"`
}

// ----- EXAMPLES -----
// passwords/email/google.com/email1@gmail.com:<path_hash>
// passwords/email/google.com/email2@gmail.com:<path_hash>
// passwords/email/tutanota.com/email1@tutanota.com:<path_hash>
// ----- EXAMPLES ------
// The above are examples of how the strings are formatted within the file, which we will use to populate the map
// The keys of the map will be the hash, and the values will be the filepath for the account file
type TreeDataBase map[string]string
