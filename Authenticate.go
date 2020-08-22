package go_email

import (
	"encoding/json"
	"errors"
	"github.com/AadityaChaudhary/go-oauthdialog"
	"github.com/AadityaChaudhary/go-sasl"
	"github.com/emersion/go-imap/client"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func authenticate(c *client.Client, cfg *oauth2.Config, username string) (sasl.Client, error) {


	spt,err := c.SupportAuth(sasl.Xoauth2)
	if err != nil {
		return nil, err
	}
	if  !spt {
		return nil, errors.New("XOAUTH2 not supported by the server")
	}

	var token *oauth2.Token

	if data ,err := ioutil.ReadFile("token.json"); err != nil {

		// Ask for the user to login with his Google account

		code, err := oauthdialog.Open(cfg)

		if err != nil {
			log.Println("opening error")
			return nil, err
		}

		// Get a token from the returned code
		// This token can be saved in a secure store to be reused later
		token, err := cfg.Exchange(oauth2.NoContext, code, oauth2.AccessTypeOffline)
		file, err := json.Marshal(token)

		_ = ioutil.WriteFile("token.json", file, os.ModePerm)


		if err != nil {
			log.Println("getting token error")
			return nil, err
		}
	} else {
		log.Println("file found")
		err := json.Unmarshal(data, &token)
		if err != nil {
			log.Fatal(err)
		}

	}


	// Login to the IMAP server with XOAUTH2
	saslClient := sasl.NewXoauth2Client(username, token.AccessToken)

	//log.Println("new xoauth client error")
	return saslClient, nil
}


func(ec *EmailClient) GetToken(code string, config *oauth2.Config) error {
	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return err
	}
	ec.Token = token
	return nil
}

func(ec *EmailClient) Authenticate(username string) error {

	ec.Auth = sasl.NewXoauth2Client(username, ec.Token.AccessToken)

	//log.Println(token)
	err := ec.Client.Authenticate(ec.Auth)
	return err

}

func(ec *EmailClient) AllInOneAuth(username string, config oauth2.Config) (oauth2.Token, error) {
	spt,err := ec.Client.SupportAuth(sasl.Xoauth2)
	if err != nil {
		return oauth2.Token{}, err
	}
	if  !spt {
		return oauth2.Token{}, errors.New("XOAUTH2 not supported by the server")
	}
	code, err := oauthdialog.Open(&config)
	if err != nil {
		log.Fatal(err)
	}
	token, err := config.Exchange(oauth2.NoContext, code, oauth2.AccessTypeOffline)
	if err != nil {
		return oauth2.Token{}, err
	}
	ec.Token = token
	//saslClient := sasl.NewXoauth2Client(username, token.AccessToken)
	//err = ec.Client.Authenticate(saslClient)
	//if err != nil {
	//	return *token,err
	//}
	return *token, err


}

func(ec *EmailClient) SetTokenVals(accessToken, expiry, refreshToken string, tokenType string) error {

	ec.Token = &oauth2.Token{}
	ec.Token.AccessToken = accessToken
	ec.Token.RefreshToken = refreshToken
	ec.Token.TokenType = tokenType

	var err error
	if ec.Token.Expiry, err = time.Parse(time.RFC3339, expiry); err != nil {
		return err
	}
	return nil
}





