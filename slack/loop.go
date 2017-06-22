package slack

import (
	"log"
	"strings"
)

// Loop ...
func Loop(conn *Conn, cmds []Command) {
	msgs := make(chan Message)
	go commandProcessor(conn, msgs, cmds)

	for {
		m, err := conn.Get()
		if err != nil {
			panic(err)
		}

		if m.Type != "message" {
			continue
		}

		if m.Text == "" {
			continue
		}

		log.Printf(">>> [%s:@%s] %s", m.Type, m.UserDetail.User.Name, m.Text)
		m.SelfID = "<@" + conn.Self.ID + ">"

		if strings.HasPrefix(m.Channel, "D") {
			m.Text = m.SelfID + " " + m.Text // hack
		}

		words := strings.Split(m.Text, " ")
		if words[0] == m.SelfID {
			m.Respond = true
			m.Text = strings.Join(words[1:], " ")
		}

		msgs <- m
	}
}
