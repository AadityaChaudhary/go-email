package go_email

import (
	"golang.org/x/oauth2"
)

type EmailConfig struct {
	Smtp 		string			`json:"smtp"`
	Imap 		string			`json:"imap"`
	AuthURI		string			`json:"authURI"`
	ExchangeURI	string			`json:"exchangeURI"`
	Config 		oauth2.Config	`json:"config"`

}
