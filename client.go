package irc

import "net"
import "fmt"
import "bufio"
import "crypto/tls"

type Client struct {
	socket net.Conn

	Host     string
	Port     int
	Nickname string
	Ident    string
	Realname string
	Secure   bool

	Handler EventHandler
}

func (c *Client) Write(s string) error {
	_, err := c.socket.Write([]byte(s + "\r\n"))

	return err
}

func (c *Client) Connect() error {
	var socket net.Conn
	var err error

	if c.Secure {
		socket, err = tls.Dial("tcp", fmt.Sprintf("%s:%v", c.Host, c.Port), &tls.Config{})
	} else {
		socket, err = net.Dial("tcp", fmt.Sprintf("%s:%v", c.Host, c.Port))
	}

	if err != nil {
		return err
	}

	c.socket = socket

	c.Write("NICK " + c.Nickname)
	c.Write("USER " + c.Ident + " 0 * :" + c.Realname)

	if err = c.readPump(); err != nil {
		return err
	}

	return nil
}

func (c *Client) readPump() error {
	reader := bufio.NewReader(c.socket)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			return err
		}

		line = line[0 : len(line)-2]

		message := Message{raw: line}
		message.parse()

		c.Handler.trigger(c, &message)
	}
}
