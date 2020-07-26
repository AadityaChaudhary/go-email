package go_email

import (
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Envelope struct {
	Sent 	time.Time
	Subject	string

	To 		[]mail.Address
	From 	[]mail.Address
	Cc 		[]mail.Address
	Bcc 	[]mail.Address
}



type Message struct {
	Parts 		[]MessagePart
	Envelope 	Envelope
}

type MessagePart struct {
	Name 		string
	PartType 	string
	Part 		[]byte
}

func(m *Message) Print() {
	log.Println("Date: ", m.Envelope.Sent)
	log.Println("Subject: ", m.Envelope.Subject)
	log.Println("To: ", m.Envelope.To)
	log.Println("From: ", m.Envelope.From)
	log.Println("CC: ", m.Envelope.Cc)
	log.Println("BCC: ", m.Envelope.Bcc)

	for _, part := range m.Parts {
		log.Println("---------------------")
		log.Println("Type: ", part.PartType)
		log.Println("Name: ", part.Name)
		log.Println("data: ", string(part.Part))
	}
}

func(ec *EmailClient) GetMailBoxes() ([]imap.MailboxInfo, error) {
	var mboxs []imap.MailboxInfo
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error)

	go func() {
		done <- ec.Client.List("","*",mailboxes)
	}()

	log.Println("pulled mailboxes")
	for m := range mailboxes {
		mboxs = append(mboxs, *m)
	}
	if err := <- done; err != nil {
		return nil, err
	}
	return mboxs, nil
}

//func(ec *EmailClient) GetMail(mName string, to uint32, from uint32) ([]imap.Message, error) {
//	//getting the inbox mailbox
//	mbox, err := ec.Client.Select(mName, true) //readonly on true for safety during testing lol
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	seqset := new(imap.SeqSet)
//	seqset.AddRange(from,to)
//
//	messages := make(chan *imap.Message,10)
//
//	done := make(chan error,1)
//
//	go func() {
//		done <- ec.Client.Fetch(seqset,[]imap.FetchItem{imap.FetchBody}, messages)
//	}()
//	x := 0
//	for msg := range messages {
//
//
//	}
//	log.Println(x)
//
//	if err:= <- done; err != nil {
//		log.Fatal(err)
//	}
//
//	ec.Client.Subscribe("")
//
//}

func(ec *EmailClient) SelectMailBox(mName string) (error) {
	_, err :=  ec.Client.Select(mName, true) //readonly on true for safety during testing lol
	return err

}

func(ec *EmailClient) GetEnvelopes(from uint32, to uint32) []imap.Envelope {

	if from > ec.Client.Mailbox().Messages || to > ec.Client.Mailbox().Messages {
		log.Fatal("ruh roh")
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from,to)

	messages := make(chan *imap.Message,10)

	done := make(chan error,1)

	go func() {
		done <- ec.Client.Fetch(seqset,[]imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	var envelopes []imap.Envelope

	for  msg := range messages {
		envelopes = append(envelopes, *msg.Envelope)
	}

	return envelopes


}

func(ec *EmailClient) GetBody(uids []uint32) (imap.Message, imap.BodySectionName, error) {
	seqset := new(imap.SeqSet)
	seqset.AddNum(uids[0])

	messages := make(chan *imap.Message,10)

	done := make(chan error,1)

	var section imap.BodySectionName

	go func() {
		done <- ec.Client.Fetch(seqset,[]imap.FetchItem{section.FetchItem()}, messages)
	}()

	for  msg := range messages {
		return *msg,section, nil
	}

	return imap.Message{},imap.BodySectionName{}, errors.New("message not found")



}

func(ec *EmailClient) GetLast(amount uint32) (uint32, uint32) {
	return ec.Client.Mailbox().Messages - amount, ec.Client.Mailbox().Messages  //from, to
}

func(ec *EmailClient) GetEnvelopesFromArr(msgs []uint32) []imap.Envelope {

	seqset := new(imap.SeqSet)
	for _, val := range msgs {
		seqset.AddNum(val)
	}

	messages := make(chan *imap.Message,10)

	done := make(chan error,1)

	go func() {
		done <- ec.Client.Fetch(seqset,[]imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}, messages)
	}()

	var envelopes []imap.Envelope

	for  msg := range messages {
		envelopes = append(envelopes, *msg.Envelope)
	}

	return envelopes


}

func(ec *EmailClient) ParseMessage(msg imap.Message, section imap.BodySectionName) (Message, error) {
	r := msg.GetBody(&section)
	if r == nil {
		log.Fatal("server didnt return message body")
	}

	var message Message
	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal(err)
	}

	header := mr.Header
	if date, err := header.Date(); err == nil {
		log.Println("date: ", date)
		message.Envelope.Sent = date
	}
	if from, err := header.AddressList("From"); err == nil {
		log.Println("From:", from)
		for _, add := range from {
			message.Envelope.From = append(message.Envelope.From, *add)
		}

	}
	if to, err := header.AddressList("To"); err == nil {
		log.Println("To:", to)
		for _, add := range to {
			message.Envelope.To = append(message.Envelope.To, *add)
		}
	}
	if cc, err := header.AddressList("Cc"); err == nil {
		log.Println("Cc:", cc)
		for _, add := range cc {
			message.Envelope.Cc = append(message.Envelope.Cc, *add)
		}
	}
	if bcc, err := header.AddressList("Bcc"); err == nil {
		log.Println("Bcc:", bcc)
		for _, add := range bcc {
			message.Envelope.Bcc = append(message.Envelope.Bcc, *add)
		}
	}
	if subject, err := header.Subject(); err == nil {
		log.Println("Subject:", subject)
		message.Envelope.Subject = subject
	}

	for {
		p, err := mr.NextPart()

		if err == io.EOF {
			break
		} else if err != nil {
			return Message{}, err
		}
		var part MessagePart
		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			b, _ := ioutil.ReadAll(p.Body)
			//log.Println("Got text: ", string(b))
			if strings.HasPrefix(string(b),"<") {
				part.PartType = "html"
			} else {
				part.PartType = "raw text"
			}
			part.Name = "text"
			part.Part = b
		case *mail.AttachmentHeader:
			filename, _ := h.Filename()
			//log.Println("Got Attachment: ", filename)
			b, _ := ioutil.ReadAll(p.Body)
			part.PartType = "attachment"
			part.Name = filename
			part.Part = b

		}

		message.Parts = append(message.Parts, part)

	}

	return message, nil
}


