# Password Manager in Go

## Steps to Encrypt and Store Passwords

- Given a master password from the user, we hash it using bcrypt, and then sha256, which we will then use as an encryption key
- At encryption time we will generate a nonce (a pseudo-random number), which will be prepended to the (encrypted) password, which would have been encrypted by the encryption key (master password).

## Password Database Structure

```json
"email": {
    "google.com": {
        "Username": {
            "Password": "password",
            "Email": "email@gmail.com",
            "Username": "username",
            "Domain": "google.com",
            "RecoveryCodes": [
                "324345345",
                "342342343"
            ]
        },
        "Username2": {
            "Password": "password2",
            "Email": "email2@gmail.com",
            "Username": "username2",
            "Domain": "google.com",
            "RecoveryCodes": [
                "dsfsdfdsf"
            ]
        }
    },
    "tutanota.com": {
        "Password3": "password3",
        "Email": "email3@tutanota.com",
        "Username": "username3",
        "Domain": "tutanota.com",
        "RecoveryCodes": [
            "asdasd",
            "325234",
            "desdf3",
            "asd223"
        ]
    },
    "protonMail": {
        "Password4": "password4",
        "Email": "email4@proton.me",
        "Username": "username4",
        "Domain": "proton.me",
        "RecoveryCodes": [
            "asdasd",
            "325234",
            "desdf3",
            "asd223"
        ],
        "extraData": {
            "Something": "something"
        }
    }
}
```
Either the way above, or the way below, but I havent decided yet, maybe have a configuration to chose either
```
DATABASE/
|
|_ <hash_for_string_"email">/
    |_ <hash_for_string_"google.com">/
        |_ <hash_for_username>.pw
        |_ <hash_for_username2>.pw
    |_ <hash_for_string_"tutanota.com">.pw
    |_ <hash_for_string_"protonMail">.pw
```
I want to query account data with the syntax show below, and regardless of how data is stored on the disk, I want to parse it in a json format
```
passman find email
> google.com/username
> google.com/username2

passman show "google.com/username"
>  "Password": "password",
>  "Email": "email@gmail.com",
>  "Username": "username",
>  "Domain": "google.com",
>  "RecoveryCodes": [
>    "324345345",
>    "342342343"
>  ]

passman show google.com/username:password
> password

passman show google.com/username2:password
> password2
```
