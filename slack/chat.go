// Package slack provides an interface for working with slack.com
package slack

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"regexp"
)

type cart struct {
	plugin  *plugin.Plugin
	regpatt string
	help    string
}

// ChatLoop enters a hard loop that reads off messages and processes them.
func ChatLoop(conn Conn) {
	carts := loadCarts()
	if os.Getenv("DEBUG") == "1" {
		for _, c := range carts {
			log.Printf("DBG cart: %+v\n", c)
		}
	}

	for {
		msg, err := conn.Get()
		if err != nil {
			log.Fatal(err)
		}

		if !msg.Respond {
			continue
		}

		log.Printf(">>> %s %s: %s (cmd=%s)\n", msg.Channel, msg.User, msg.Text, msg.Command)

		// built-ins
		if msg.User == "U52JX5HPE" { // TODO auth system
			switch {
			case msg.Command == "reload":
				carts = loadCarts()
			case msg.Command == "dump":
				dump(msg, conn)
			case msg.Command == "help":
				help(msg, conn, carts)
			}
		}

		// match input
		for _, cart := range carts {
			re := regexp.MustCompile(cart.regpatt)
			m := re.FindStringSubmatch(msg.Full)

			if os.Getenv("DEBUG") == "1" {
				log.Printf("DBG re=%s line=%s", re, msg.Full)
			}

			resp, err := cart.plugin.Lookup("Respond")
			if err != nil {
				log.Printf("ERR %s\n", err)
				continue
			}

			if len(m) > 0 {
				go resp.(func(Message, Conn, []string))(msg, conn, m)
			}
		}
	}
}

func dump(m Message, c Conn) {
	fmt.Printf("MSG:\n%+v\n\nCONN:\n%+v\n\n", m, c)
}

func loadCarts() (carts []cart) {
	// TODO reloading doesn't work as plugin.Open always returns the
	//      same *Plugin
	dir := "plugins" // TODO config
	// TODO shitty path handling
	cartfiles, err := filepath.Glob(dir + "/*.so")
	if err != nil {
		log.Printf("ERR %s\n", err)
		return
	}
	var c cart
	for _, cartfile := range cartfiles {
		log.Printf("-i- loadcart: %s\n", cartfile)
		p, err := plugin.Open(cartfile)
		if err != nil {
			log.Printf("ERR %s\n", err)
			continue
		}
		c = cart{plugin: p}
		if register(&c) {
			carts = append(carts, c)
		}
	}

	return
}

func register(c *cart) bool {
	r, err := c.plugin.Lookup("Register")
	if err != nil {
		log.Printf("ERR %s\n", err)
		return false
	}
	c.regpatt, c.help = r.(func() (string, string))()
	return true
}

func help(m Message, c Conn, carts []cart) {
	reply := Reply{}
	reply.Channel = m.Channel
	if len(carts) < 1 {
		reply.Text = "Perdone, pero creo que le han desinformado."
		c.Send(m, reply)
		return
	}
	reply.Text = "Escoja lo que prefiera, invita la casa:\n"
	for _, cart := range carts {
		reply.Text += ":point_right: " + cart.help + "\n"
	}
	c.Send(m, reply)
}
