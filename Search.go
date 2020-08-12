package go_email

import (
	"github.com/emersion/go-imap"
	"time"
)

type Crit struct {
	Since 		time.Time
	Before 		time.Time
	SentSince 	time.Time
	SentBefore 	time.Time

	From 		[]string
	To 			[]string
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
	for _, val := range c.To {
		criteria.Header.Add("to", val)
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

	criteria.Before = c.Before
	criteria.Since = c.Since
	criteria.SentBefore = c.SentBefore
	criteria.SentSince = c.SentSince

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













