package go_email

import (
	"golang.org/x/oauth2"
)

type EmailConfig struct {
	Defaults 	map[string]string 	`json:"defaults"`
	Smtp 		string				`json:"smtp"`
	Imap 		string				`json:"imap"`
	AuthURI		string				`json:"authURI"`
	ExchangeURI	string				`json:"exchangeURI"`
	Config 		oauth2.Config		`json:"config"`

}

type Defaults struct {
	Inbox 		string 			`json:"inbox"`
	Trash 		string 			`json:"trash"`
	Spam 		string 			`json:"spam"`

}
