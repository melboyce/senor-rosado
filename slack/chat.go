// Package slack provides an interface for working with slack.com
package slack

import (
	"fmt"
	"log"
	"path/filepath"
	"plugin"
	"regexp"
)

type cart struct {
	plugin  *plugin.Plugin
	regpatt string
}

// ChatLoop enters a hard loop that reads off messages and processes them.
func ChatLoop(conn Conn) {
	carts := loadCarts()
	for _, c := range carts {
		log.Printf("DBG cart: %+v\n", c)
	}

	for {
		msg, err := conn.Get()
		if err != nil {
			log.Fatal(err)
		}

		if !msg.Respond {
			// Message.Respond: if true, message is targetted at bot
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
			}
		}

		// match input
		for _, cart := range carts {
			re := regexp.MustCompile(cart.regpatt)
			m := re.FindStringSubmatch(msg.Full)
			resp, err := cart.plugin.Lookup("Respond")
			if err != nil {
				log.Printf("ERR %s\n", err)
				continue
			}
			if len(m) > 0 {
				resp.(func(Message, Conn, []string))(msg, conn, m)
			}
		}
	}
}

func dump(m Message, c Conn) {
	fmt.Printf("MSG:\n%+v\n\nCONN:\n%+v\n\n", m, c)
}

func loadCarts() (carts []cart) {
	dir := "plugins" // TODO config
	// TODO shitty path handling
	cartfiles, err := filepath.Glob(dir + "/*.so")
	if err != nil {
		log.Printf("ERR %s\n", err)
		return
	}
	carts = make([]cart, len(cartfiles), len(cartfiles))
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

	return carts[1:]
}

func register(c *cart) bool {
	r, err := c.plugin.Lookup("Register")
	if err != nil {
		log.Printf("ERR %s\n", err)
		return false
	}
	c.regpatt = r.(func() string)()
	return true
}
