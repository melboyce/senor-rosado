// Package slack provides an interface for working with slack.com
package slack

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"time"

	"encoding/json"
	"net/http"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

// A Conn represents a slack connection.
type Conn struct {
	Ok   bool   `json:"ok"`
	URL  string `json:"url"`
	Team struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
	Self struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"self"`
	Error string `json:"error"`

	Sock *websocket.Conn
}

// A Message is a slack RTM message object with some meta.
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`

	Respond     bool
	Target      string
	Full        string
	Command     string
	Subcommand  string
	Tail        string
	SubTail     string
	ReplyToUser bool
}

// A Reply is another name for a Message
type Reply Message

// CartLibrary is a slice of cartirdges and some meta
type CartLibrary struct {
	Carts []Cart
}

// Cart is a single plugin cartridge
type Cart struct {
	Plugin  *plugin.Plugin
	Regpatt string
	Help    string
}

// see: Conn.Send()
var counter uint64

var httpClient = &http.Client{Timeout: 10 * time.Second}

var slackAPIURL = "https://slack.com/api/rtm.connect?token=%s"

// Connect to slack and return a useful struct or an error.
func Connect(token string) (slack Conn, err error) {
	url := fmt.Sprintf(slackAPIURL, token)
	r, err := httpClient.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		err = fmt.Errorf("Error: status code: %d", r.StatusCode)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&slack)
	if err != nil {
		return
	}
	if !slack.Ok {
		err = fmt.Errorf("Error: slack: %s", slack.Error)
		return
	}

	// attach a websocket for RTM comms to the slack struct
	slack.Sock, err = websocket.Dial(slack.URL, "", "https://api.slack.com/")
	if err != nil {
		return
	}

	return
}

// Get pulls a message out of the RTM queue and returns it as a Message struct.
func (s *Conn) Get() (m Message, err error) {
	err = websocket.JSON.Receive(s.Sock, &m)
	if err != nil {
		return
	}

	if m.Type != "message" {
		return
	}

	// simple tokenization
	words := strings.Split(m.Text, " ")

	// slack.ChatLoop uses this to ignore messages not targetted at the bot
	if words[0] == "<@"+s.Self.ID+">" {
		m.Respond = true
		m.Target = words[0] // convenience
	}
	if strings.HasPrefix(m.Channel, "D") {
		m.Respond = true
		m.Target = s.Self.ID
		words = append([]string{s.Self.ID}, words...)
	}

	// these are conveniences used when authoring commands
	n := len(words)
	if n > 1 {
		m.Command = words[1]
		m.Full = strings.Join(words[1:], " ")
	}
	if n > 2 {
		m.Subcommand = words[2]
		m.Tail = strings.Join(words[2:], " ")
	}
	if n > 3 {
		m.SubTail = strings.Join(words[3:], "")
	}

	return
}

// Send pushes a Reply struct onto the RTM queue
func (s *Conn) Send(m Message, r Reply) error {
	r.ID = atomic.AddUint64(&counter, 1)
	r.Type = "message" // TODO this will bite me later
	if m.User != "" && r.ReplyToUser {
		r.Text = "<@" + m.User + "> " + r.Text
	}
	log.Printf("<<< %s\n", r.Text)
	return websocket.JSON.Send(s.Sock, &r)
}

// Load populates CartLibrary.Carts with plugins
func (l *CartLibrary) Load() (err error) {
	dir := os.Getenv("SR_PLUGIN_DIR")
	if dir == "" {
		dir = "plugins"
	}

	glob, err := filepath.Glob(filepath.Join(dir, "*.so"))
	if err != nil {
		return
	}

	var c Cart
	for _, fp := range glob {
		log.Printf("INF loading cart: %s\n", fp)
		p, err := plugin.Open(fp)
		if err != nil {
			fmt.Printf("ERR %s\n", err)
			continue
		}

		c = Cart{Plugin: p}
		if register(&c) {
			l.Carts = append(l.Carts, c)
		}
	}

	return
}

func register(c *Cart) bool {
	r, err := c.Plugin.Lookup("Register")
	if err != nil {
		log.Printf("ERR %s\n", err)
		return false
	}
	c.Regpatt, c.Help = r.(func() (string, string))()
	return true
}

// GetJSON is a helper for unmarshalling JSON responses
func GetJSON(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		err = fmt.Errorf("ERR non-200 response: \"%d\" url=%s", r.StatusCode, url)
		return err
	}

	return json.NewDecoder(r.Body).Decode(target)
}

// PanicSuppress is a helper for plugins to blind-recover from panic()
func PanicSuppress() {
	if r := recover(); r != nil {
		fmt.Println(".!. PANIC SUPRESSED:", r)
	}
}
