package slack

import (
	"log"
	"regexp"
)

// Command ...
type Command struct {
	Register func() (string, string)
	Respond  func(Reply, chan Reply)
	Re       *regexp.Regexp
}

// TODO plugins
func commandProcessor(conn *Conn, msgs chan Message, cmds []Command) {
	help := "help:\n"
	for i, cmd := range cmds {
		p, h := cmd.Register()
		help += h + "\n"
		cmds[i].Re = regexp.MustCompile(p)
	}

	replies := make(chan Reply)
	var m Message
	var r Reply
	var err error
	for {
		select {
		case m = <-msgs:
			matchCommands(m, cmds, replies)
		case r = <-replies:
			if err = conn.Send(r); err != nil {
				log.Printf("!!! COMM PROC: %s", err)
			}
		}
	}
}

func matchCommands(m Message, cmds []Command, replies chan Reply) {
	if !m.Respond {
		return
	}
	r := GetReply(m)
	for _, cmd := range cmds {
		if cmd.Re == nil {
			continue
		}
		matches := cmd.Re.FindAllStringSubmatch(m.Text, -1)
		if len(matches) > 0 && len(matches[0]) > 0 {
			r.Matches = matches
			go cmd.Respond(r, replies)
		}
	}
}
