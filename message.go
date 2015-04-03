package irc

import "strings"

type Message struct {
	raw string

	Prefix  string
	Command string
	Params  []string
}

func (m *Message) parse() {
	prefixEnd := -1
	trailingStart := len(m.raw)
	trailing := ""

	if strings.HasPrefix(m.raw, ":") {
		prefixEnd = strings.Index(m.raw, " ")
		m.Prefix = m.raw[1:prefixEnd]
	}

	commandAndParameters := ""
	trailingStart = strings.Index(m.raw, " :")
	if trailingStart >= 0 {
		trailing = m.raw[trailingStart+2:]
		commandAndParameters = m.raw[prefixEnd+1 : trailingStart]
	} else {
		commandAndParameters = m.raw[prefixEnd+1:]
	}

	commandAndParametersSegments := strings.Split(commandAndParameters, " ")

	m.Command = commandAndParametersSegments[0]

	if len(commandAndParametersSegments) > 1 {
		m.Params = commandAndParametersSegments[1:]
	}

	if trailingStart >= 0 {
		m.Params = append(m.Params, trailing)
	}
}
