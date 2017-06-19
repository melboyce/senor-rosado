package slack

import (
	"fmt"
	"log"
	"strings"
)

// Loop ...
func Loop(conn *Conn) (err error) {
	for {
		m, err := conn.Get()
		if err != nil {
			return
		}

		if m.Type == "message" {
			continue
		}

		log.Printf(">>> [%s] %s", m.Type, m.Text)

		// NOTE mutates Message
		if err = processMessage(&m); err != nil {
			log.Printf("/!\\ LOOP PROC: %s", err)
		}
	}
}

func processMessage(conn Conn, m *Message) (err error) {
	if m.Text == "" {
		err = fmt.Errorf("-w- LOOP PROC: empty message")
		return
	}

	words := strings.Split(m.Text, " ")
	if words[0] == "<@"+conn.Self.ID+">" {
		m.Respond = true
		m.Text = strings.Join(words[1:], " ")
	}
	if strings.HasPrefix(m.Channel, "D") {
		m.Respond = true
		selfID := "<@" + conn.Self.ID + ">"
		words = append([]string{selfID}, words...) // hack!
	}

	n := len(words)
	if n > 0 {
		m.Cmd = words[1]
	}
	if n > 1 {
		m.Args = words[2:]
	}
}
