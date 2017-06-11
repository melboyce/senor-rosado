// Package slack provides an interface for working with slack.com
package slack

import (
	"fmt"
	"log"
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

// see: Conn.Send()
var counter uint64

var httpClient = &http.Client{Timeout: 10 * time.Second}

// Connect to slack and return a useful struct or an error.
func Connect(token string) (slack Conn, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.connect?token=%s", token)
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
	}

	// these are conveniences used when authoring commands
	m.Target = words[0] // the user that was targetted
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
	log.Printf("<<< %+v\n", r)
	return websocket.JSON.Send(s.Sock, &r)
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

// PanicSuppress ...
func PanicSuppress() {
	if r := recover(); r != nil {
		fmt.Println(".!. PANIC SUPRESSED:", r)
	}
}
