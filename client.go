package irc

import "net"
import "fmt"
import "bufio"
import "crypto/tls"
import "time"

type Client struct {
	conn net.Conn

	Host     string
	Port     int
	Nickname string
	Ident    string
	Realname string
	Secure   bool
	Password string

	Handler EventHandler
}

func (c *Client) Write(s string) error {
	_, err := c.conn.Write([]byte(s + "\r\n"))

	return err
}

func (c *Client) Writef(format string, a ...interface{}) error {
	return c.Write(fmt.Sprintf(format, a...))
}

func (c *Client) Connect() error {
	var conn net.Conn
	var err error

	if c.Secure {
		conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%v", c.Host, c.Port), &tls.Config{})
	} else {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%v", c.Host, c.Port))
	}

	if err != nil {
		return err
	}

	c.conn = conn

	go c.setupPingLoop()
	go c.setupReadLoop()

	if len(c.Password) > 0 {
		c.Writef("PASS %s", c.Password)
	}

	c.Writef("NICK %s", c.Nickname)
	c.Writef("USER %s 0 * :%s", c.Ident, c.Realname)

	return nil
}

func (c *Client) setupPingLoop() {
	ticker := time.NewTicker(time.Minute * 1)

	for _ = range ticker.C {
		c.Writef("PING :%d", time.Now().UnixNano())
	}
}

func (c *Client) setupReadLoop() {
	reader := bufio.NewReader(c.conn)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			// todo
			panic(err)
		}

		line = line[0 : len(line)-2]

		message := Message{raw: line}
		message.parse()

		c.Handler.trigger(c, &message)
	}
}
