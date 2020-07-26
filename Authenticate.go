package go_email

import (
	"encoding/json"
	"errors"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-oauthdialog"
	"github.com/emersion/go-sasl"
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
		token, err := cfg.Exchange(oauth2.NoContext, code)
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

func(ec *EmailClient) GetAuthURL() (string, string, error) {
	ec.Config.Config.RedirectURL = ec.Config.AuthURI
	spt,err := ec.Client.SupportAuth(sasl.Xoauth2)
	if err != nil {
		return "","", err
	}
	if  !spt {
		return "", "", errors.New("XOAUTH2 not supported by the server")
	}
	state, err := oauthdialog.GenerateState()
	if err != nil {
		log.Fatal(err)
	}
	url := ec.Config.Config.AuthCodeURL(state)
	return url, state, nil
}

func(ec *EmailClient) GetToken(code string) error {
	ec.Config.Config.RedirectURL = ec.Config.ExchangeURI
	token, err := ec.Config.Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return err
	}
	ec.Token = token
	return nil
}

func(ec *EmailClient) Authenticate(username string) error {

	//log.Println(token)
	err := ec.Client.Authenticate(sasl.NewXoauth2Client(username, ec.Token.AccessToken))
	return err

}

func(ec *EmailClient) AllInOneAuth(username string) (oauth2.Token, error) {
	spt,err := ec.Client.SupportAuth(sasl.Xoauth2)
	if err != nil {
		return oauth2.Token{}, err
	}
	if  !spt {
		return oauth2.Token{}, errors.New("XOAUTH2 not supported by the server")
	}
	code, err := oauthdialog.Open(&ec.Config.Config)
	if err != nil {
		log.Fatal(err)
	}
	token, err := ec.Config.Config.Exchange(oauth2.NoContext, code)
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




