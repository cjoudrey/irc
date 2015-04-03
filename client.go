package irc

import "net"
import "fmt"
import "bufio"
import "crypto/tls"
import "time"

type Client struct {
	conn  net.Conn
	write chan string

	Host     string
	Port     int
	Nickname string
	Ident    string
	Realname string
	Secure   bool
	Password string

	Handler EventHandler
}

func (c *Client) Write(s string) {
	c.write <- s
}

func (c *Client) Writef(format string, a ...interface{}) {
	c.Write(fmt.Sprintf(format, a...))
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
	c.write = make(chan string)

	go c.setupPingLoop()
	go c.setupReadLoop()
	go c.setupWriteLoop()

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

func (c *Client) setupWriteLoop() {
	tickets := make(chan bool, 5)
	ticker := time.Tick(1 * time.Second)

	go func() {
		for _ = range ticker {
			tickets <- true
		}
	}()

	for line := range c.write {
		<-tickets
		_, err := c.conn.Write([]byte(line + "\r\n"))

		if err != nil {
			// todo
			panic(err)
		}
	}
}
