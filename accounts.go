package main

type Account struct {
	Password      string       `json:"Password"`
	Username      string       `json:"Username"`
	Email         string       `json:"Email"`
	Service       string       `json:"Service"`
	RecoveryCodes []string     `json:"RecoveryCodes"`
	ExtraData     *interface{} `json:"ExtraData"`
}
