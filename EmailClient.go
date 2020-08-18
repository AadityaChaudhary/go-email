package go_email

import (
	"encoding/json"
	"github.com/AadityaChaudhary/go-sasl"
	"github.com/emersion/go-imap-idle"
	"github.com/emersion/go-imap/client"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
)

type EmailClient struct {
	Config 			EmailConfig
	Token			*oauth2.Token
	Client			*client.Client
	IdleClient 		*idle.Client
	EmailAddress	string
	Auth 		 	sasl.Client

}

func NewEmailClient(configsPath string, cfg string) (EmailClient,error) {

		var emailClient EmailClient
		emailClient.EmailAddress = cfg

		var configs map[string]EmailConfig
		if data, err := ioutil.ReadFile(configsPath); err != nil {
			return EmailClient{}, err
		} else {
			err := json.Unmarshal(data,&configs)
			if err != nil {
				return EmailClient{}, err
			}
		}
		emailClient.Config = configs[cfg]
		log.Println(emailClient.Config.Imap)

		return emailClient, nil
}

func(ec *EmailClient) DialTlsImap() error {
	var err error
	log.Println(ec.Config.Imap,": imap")
	ec.Client, err = client.DialTLS(ec.Config.Imap, nil)
	ec.IdleClient = idle.NewClient(ec.Client)

	if err != nil {
		return err
	}
	return nil
}

func(ec *EmailClient) LogOut() error{
	err := ec.Client.Logout()
	return err
}



