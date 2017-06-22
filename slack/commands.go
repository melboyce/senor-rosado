package slack

import (
	"log"
	"regexp"
)

// Command contains the two functions exported by a command and the compiled
// regexp used to match lines.
type Command struct {
	Register func() (string, string)
	Respond  func(Reply, chan Reply)
	Re       *regexp.Regexp
}

var helpMessage string

// commandProcessor is a goroutine started prior to the main loop.
func commandProcessor(conn *Conn, msgs chan Message, cmds []Command) {
	helpMessage = "help:\n"
	for i, cmd := range cmds {
		p, h := cmd.Register()
		helpMessage += h + "\n"
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

// matchCommands checks a Message against the regexps and potentially launches
// a command as a goroutine in response.
func matchCommands(m Message, cmds []Command, replies chan Reply) {
	if !m.Respond {
		return
	}
	for _, cmd := range cmds {
		if cmd.Re == nil {
			continue
		}
		r := GetReply(m)
		matches := cmd.Re.FindAllStringSubmatch(m.Text, -1)
		if len(matches) > 0 && len(matches[0]) > 0 {
			r.Matches = matches
			go cmd.Respond(r, replies)
		}
	}
}
