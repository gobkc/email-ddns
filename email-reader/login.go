package email_reader

import (
	"github.com/emersion/go-imap/client"
	"net"
	"time"
)

func Login(server, user, pass string) (*client.Client, error) {
	dial := new(net.Dialer)
	dial.Timeout = time.Duration(3) * time.Second
	c, err := client.DialWithDialerTLS(dial, server, nil)
	if err != nil {
		c, err = client.DialWithDialer(dial, server)
	}
	if err != nil {
		return nil, err
	}
	if err = c.Login(user, pass); err != nil {
		return nil, err
	}
	return c, nil
}
