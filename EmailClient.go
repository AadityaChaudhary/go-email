package go_email

import (
	"encoding/json"
	"github.com/AadityaChaudhary/go-sasl"
	"github.com/emersion/go-imap/client"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
)

type EmailClient struct {
	Platform 		string
	Token			*oauth2.Token
	Client			*client.Client

	EmailAddress	string
	Auth 		 	sasl.Client

}

func NewEmailClient(cfg string) (EmailClient,error) {

		var emailClient EmailClient
		emailClient.Platform = cfg

		return emailClient, nil
}

func NewConfig(configsPath string) (map[string]EmailConfig, error) {
	var configs map[string]EmailConfig
	if data, err := ioutil.ReadFile(configsPath); err != nil {
		return nil, err
	} else {
		err := json.Unmarshal(data,&configs)
		if err != nil {
			return nil, err
		} else {
			return configs, nil
		}
	}
}


func(ec *EmailClient) DialTlsImap(imap string) error {
	var err error
	log.Println(imap,": imap")
	ec.Client, err = client.DialTLS(imap, nil)

	if err != nil {
		return err
	}
	return nil
}

func(ec *EmailClient) LogOut() error{
	err := ec.Client.Logout()
	return err
}



