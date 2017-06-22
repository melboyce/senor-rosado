package slack

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/websocket"
)

// Conn is a connection to Slack. It includes a Sock for sending messages.
type Conn struct {
	OK   bool   `json:"ok"`
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

	Sock  *websocket.Conn
	token string // TODO not sure if safe
}

var slackURL = "https://slack.com/api/rtm.connect?token=%s"

// Connect to slack and return a Conn struct with a connected reply Sock.
func Connect(token string) (conn Conn, err error) {
	conn.token = url.QueryEscape(token)
	u := fmt.Sprintf(slackURL, conn.token)
	if err = connect(u, &conn); err != nil {
		return
	}
	err = attachSock(&conn)
	return
}

// Get pulls a message from Slack and sets up some meta-data so it can be used.
func (conn Conn) Get() (m Message, err error) {
	if conn.Sock == nil {
		panic("!!! CONN SOCK MISSING")
	}
	if err = websocket.JSON.Receive(conn.Sock, &m); err != nil {
		return
	}
	if err = processMessage(&conn, &m); err != nil {
		return
	}
	if m.Type == "message" {
		log.Printf(">>> [#%s @%s] %s", m.ChannelDetail.Channel.Name, m.UserDetail.User.Name, m.Text)
	}
	return
}

// Send pushes a Reply on to Slack's RTM queue.
func (conn Conn) Send(r Reply) (err error) {
	if r.Text == "" {
		err = fmt.Errorf("Reply.Text is empty")
		return
	}
	log.Printf("<<< [#%s] %s", r.ChannelName, r.Text)
	return websocket.JSON.Send(conn.Sock, &r)
}

// connect performs the actual connection to Slack.
func connect(u string, conn *Conn) (err error) {
	log.Printf("-i- CONN STRT")
	if err = GetJSON(u, &conn); err != nil {
		return
	}
	if !conn.OK {
		err = fmt.Errorf("Conn.OK is false")
	}
	return
}

// attachSock dials a websocket and puts the func in the Conn struct.
func attachSock(conn *Conn) (err error) {
	log.Printf("-i- CONN SOCK: %s", conn.URL)
	origin := "https://api.slack.com/"
	conn.Sock, err = websocket.Dial(conn.URL, "", origin)
	return
}

// processMessage adds some meta-data to a Message to make it useful.
func processMessage(conn *Conn, m *Message) (err error) {
	if m.Type != "message" {
		return
	}

	if m.UserDetail, err = GetUser(conn.token, m.User); err != nil {
		return
	}

	if m.ChannelDetail, err = GetChannel(conn.token, m.Channel); err != nil {
		return
	}

	m.SelfID = "<@" + conn.Self.ID + ">"

	if strings.HasPrefix(m.Channel, "D") {
		m.Text = m.SelfID + " " + m.Text // hack
	}
	words := strings.Split(m.Text, " ")
	if words[0] == m.SelfID {
		m.Respond = true
		m.Text = strings.Join(words[1:], " ")
	}

	return
}
