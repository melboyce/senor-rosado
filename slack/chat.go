// Package slack provides an interface for working with slack.com
package slack

import (
	"log"
	"regexp"
)

// ChatLoop enters a hard loop that reads off messages and processes them.
func ChatLoop(conn Conn) {
	lib := CartLibrary{}
	if err := lib.Load(); err != nil {
		panic(err)
	}

	for {
		msg, err := conn.Get()
		if err != nil {
			log.Fatal(err)
		}

		// TODO support for commands that check all conversation
		if !msg.Respond {
			continue
		}

		log.Printf(">>> %s %s: %s (cmd=%s)\n", msg.Channel, msg.User, msg.Text, msg.Command)

		// built-ins
		switch {
		case msg.Command == "help":
			help(msg, conn, lib)
		}

		// match input
		for _, cart := range lib.Carts {
			re := regexp.MustCompile(cart.Regpatt)
			m := re.FindStringSubmatch(msg.Full)

			if len(m) > 0 {
				resp, err := cart.Plugin.Lookup("Respond")
				if err != nil {
					log.Printf("ERR %s\n", err)
					continue
				}

				// TODO find out if calling a plguin func as a goroutine is sensible
				go resp.(func(Message, Conn, []string))(msg, conn, m)
			}
		}
	}
}

func help(m Message, c Conn, l CartLibrary) {
	reply := Reply{}
	reply.Channel = m.Channel
	if len(l.Carts) < 1 {
		reply.Text = "Perdone, pero creo que le han desinformado."
		c.Send(m, reply)
		return
	}
	reply.Text = "Escoja lo que prefiera, invita la casa:\n"
	for _, cart := range l.Carts {
		reply.Text += ":point_right: " + cart.Help + "\n"
	}
	c.Send(m, reply)
}
