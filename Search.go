package go_email

import (
	"github.com/emersion/go-imap"
	"time"
)

type Crit struct {
	Since 		string
	Before 		string
	SentSince 	string
	SentBefore 	string

	From 		[]string
	Bcc 		[]string
	Cc 			[]string

	Body 		[]string
	Text 		[]string

	WithFlags	[]string
}



func(c *Crit) ToSearchCriteria() (*imap.SearchCriteria, error) {
	criteria := imap.NewSearchCriteria()

	for _, val := range c.From {
		criteria.Header.Add("from", val)
	}
	for _, val := range c.Cc {
		criteria.Header.Add("cc", val)
	}
	for _, val := range c.Bcc {
		criteria.Header.Add("bcc", val)
	}

	criteria.Body = c.Body
	criteria.Text = c.Text
	criteria.WithFlags = c.WithFlags

	var err error

	if c.Before != "" {
		criteria.Before, err = time.Parse(time.RFC3339, c.Before)
	}
	if c.Since != "" {
		criteria.Since, err = time.Parse(time.RFC3339, c.Since)
	}
	if c.SentBefore != "" {
		criteria.SentBefore, err = time.Parse(time.RFC3339, c.SentBefore)
	}
	if c.SentSince != "" {
		criteria.SentSince, err = time.Parse(time.RFC3339, c.SentSince)
	}
	if err != nil {
		return nil, err
	}

	return criteria, nil
}

func(ec *EmailClient) GetByUID(uid uint32) ([]uint32, error) {

	crit := imap.NewSearchCriteria()

	crit.Uid = new(imap.SeqSet)
	crit.Uid.AddNum(uid)

	resp, err := ec.Client.Search(crit)

	return resp, err

}

func(ec *EmailClient) GetByMessageId(messageId string) ([]uint32, error) {
	crit := imap.NewSearchCriteria()

	crit.Header.Set("Message-ID", messageId)

	resp, err := ec.Client.Search(crit)

	return resp, err
}

func(ec *EmailClient) SearchByCrit(crit Crit) ([]uint32, error) {
	criteria, err := crit.ToSearchCriteria()
	if err != nil {
		return nil, err
	}
	resp, err := ec.Client.Search(criteria)
	return resp, err
}













