package go_email

import (
	"errors"
	"github.com/emersion/go-smtp"
	"log"
	"strings"
)

//func(ec *EmailClient) SendMail( to []string, bcc []string, cc []string, subject string, msgBody string) error {
//
//	c, err := smtp.Dial(ec.Config.Smtp)
//
//	if err != nil {
//		return err
//	}
//
//	var msg string
//	//from
//	if err := c.Mail(ec.EmailAddress, nil); err != nil {
//		return err
//	}
//	msg  = msg + "From: " + ec.EmailAddress + "\r\n"
//
//	//to
//	for index, email := range to {
//		if err := c.Rcpt(email); err != nil {
//			return err
//		}
//
//		if index == 0 {
//			msg = msg + "To: " + email
//		} else {
//			msg = msg + "," + email
//		}
//	}
//
//	msg = msg + "\r\n"
//
//	//cc
//	for index, email := range cc {
//		if err := c.Rcpt(email); err != nil {
//			return err
//		}
//		if index == 0 {
//			msg = msg + "Cc: " + email
//		} else {
//			msg = msg + "," + email
//		}
//	}
//	msg = msg + "\r\n"
//
//	//bcc
//	for index, email := range bcc {
//		if err := c.Rcpt(email); err != nil {
//			return err
//		}
//		if index == 0 {
//			msg = msg + "Bcc: " + email
//		} else {
//			msg = msg + "," + email
//		}
//	}
//	msg = msg + "\r\n"
//
//	msg = msg + "Subject: " + subject + "\r\n" + "\r\n"
//
//	// Send the email body.
//	wc, err := c.Data()
//	if err != nil {
//		return err
//	}
//
//
//
//	_, err = fmt.Fprintf(wc,msg)
//	if err != nil {
//		return err
//	}
//
//
//
//	err = wc.Close()
//	if err != nil {
//		return err
//	}
//
//	// Send the QUIT command and close the connection.
//	err = c.Quit()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func(ec *EmailClient) SendMail(smtpAddress string, to []string, bcc []string, cc []string, replyTo string, subject string, msgBody string) error {


	var msg string
	//from
	msg  = msg + "From: " + ec.EmailAddress + "\r\n"

	Recp := []string{}
	//to
	for index, email := range to {

		Recp = append(Recp,email)

		if index == 0 {
			msg = msg + "To: " + email
		} else {
			msg = msg + "," + email
		}
	}
	if len(to) != 0 {
		msg = msg + "\r\n"
	} else {
		return errors.New("no <To> value specified")
	}


	//cc
	for index, email := range cc {
		Recp = append(Recp,email)

		if index == 0 {
			msg = msg + "Cc: " + email
		} else {
			msg = msg + "," + email
		}
	}
	if len(cc) != 0 {
		msg = msg + "\r\n"
	}

	//bcc
	for index, email := range bcc {
		Recp = append(Recp,email)

		if index == 0 {
			msg = msg + "Bcc: " + email
		} else {
			msg = msg + "," + email
		}
	}
	if len(bcc) != 0 {
		msg = msg + "\r\n"
	}

	if !strings.EqualFold(replyTo,"") {
		msg = msg + "In-Reply-To: " + replyTo + "\r\n"
	}

	if !strings.EqualFold(subject,"") {
		msg = msg + "Subject: " + subject + "\r\n"
	}

	// Send the email body.
	msg = msg + "\r\n" + msgBody + "\r\n"

	log.Println(msg)
	msgSend := strings.NewReader(msg)

	err := smtp.SendMail(smtpAddress,ec.Auth, ec.EmailAddress, Recp, msgSend)

	return err
}

