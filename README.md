# irc

`irc` is a extendable IRC client library written in go.

An example usage of this library can be found at: https://github.com/cjoudrey/go-irc-bot.

**This is still work in progress and should probably not be used in production. This was really just written as a learning exercise.**

## Usage

```go
package main

import "github.com/cjoudrey/irc"

func main() {
  handler := *irc.NewEventHandler()

  client := irc.Client{
    Host:     "irc.freenode.net",
    Port:     "6697",
    Nickname: "cjoudrey",
    Ident:    "cjoudrey",
    Realname: "Christian Joudrey",
    Secure:   true,
    Handler:  handler,
  }

  handler.On("001", func(c *irc.Client, m *irc.Message) {
    c.Write("JOIN #go-nuts")
  })

  client.Connect()
}
```
