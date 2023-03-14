package email_reader

import (
	"errors"
	"fmt"
	"github.com/emersion/go-imap"
	emailClientId "github.com/emersion/go-imap-id"
	"github.com/emersion/go-message/mail"
	"regexp"
	"time"
)

var ErrNotFound = errors.New("ErrNotFound")
var DomainFilter = `%s-IP changed`

func GetIp(server, user, pass, domain string) (ip string, err error) {
	client, err := Login(server, user, pass)
	if err != nil {
		return
	}
	defer func() {
		if client != nil {
			client.Close()
		}
	}()
	idClient := emailClientId.NewClient(client)
	idClient.ID(
		emailClientId.ID{
			emailClientId.FieldName:    "IMAPClient",
			emailClientId.FieldVersion: "2.1.0",
		},
	)
	boxes := make(chan *imap.MailboxInfo, 10)
	boxesDone := make(chan error, 1)
	go func() {
		boxesDone <- client.List("", "*", boxes)
	}()
	for box := range boxes {
		mbox, err := client.Select(box.Name, false)
		if err != nil {
			return "", err
		}
		if mbox.Messages == 0 {
			continue
		}
		criteria := imap.NewSearchCriteria()
		criteria.Since = time.Now().Add(-365 * time.Hour * 24)
		ids, err := client.UidSearch(criteria)
		if err != nil {
			continue
		}
		if len(ids) == 0 {
			continue
		}
		seqSet := new(imap.SeqSet)
		seqSet.AddNum(ids...)
		sect := &imap.BodySectionName{Peek: true}
		messages := make(chan *imap.Message, 100)
		messageDone := make(chan error, 1)
		go func() {
			messageDone <- client.UidFetch(seqSet, []imap.FetchItem{sect.FetchItem()}, messages)
		}()
		for msg := range messages {
			r := msg.GetBody(sect)
			mr, err := mail.CreateReader(r)
			if err != nil {
				return "", err
			}
			header := mr.Header
			ips := header.Get("X-Originating-Ip")
			subject, _ := header.Subject()
			if filter := fmt.Sprintf(DomainFilter, domain); filter == subject {
				p := regexp.MustCompile(`[0-9\.]+`)
				ipList := p.FindAllString(ips, -1)
				if len(ipList) > 0 {
					ip = ipList[0]
				}
				return ip, err
			}
		}
	}
	return "", ErrNotFound
}
