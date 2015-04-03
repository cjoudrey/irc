package irc

type EventHandler struct {
	handlers map[string][]func(client *Client, message *Message)
}

func NewEventHandler() *EventHandler {
	return &EventHandler{handlers: map[string][]func(client *Client, message *Message){}}
}

func (h *EventHandler) On(command string, callback func(client *Client, message *Message)) {
	h.handlers[command] = append(h.handlers[command], callback)
}

func (h *EventHandler) trigger(client *Client, message *Message) {
	for _, listener := range h.handlers[message.Command] {
		listener(client, message)
	}

	for _, listener := range h.handlers["*"] {
		listener(client, message)
	}
}
