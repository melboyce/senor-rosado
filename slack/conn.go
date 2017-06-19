package slack

import (
	"fmt"
	"log"
	"net/url"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

// see Conn.Send()
var counter uint64

// Conn is a Slack rtm.connect response.
type Conn struct {
	Ok   bool   `json:"ok"`
	URL  string `json:"url"`
	Team struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		Domain         string `json:"domain"`
		EnterpriseID   string `json:"enterprise_id"`
		EnterpriseName string `json:"enterprise_name"`
	} `json:"team"`
	Self struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"self"`

	Origin   string
	Sock     *websocket.Conn
	Commands []Command
}

// Connect to Slack and return a useful struct or an error
func Connect(token string) (conn Conn, err error) {
	log.Printf("-i- CONN STRT")
	token = url.QueryEscape(token)
	u := fmt.Sprintf(slackURL, token)
	err = GetJSON(u, &conn)
	if err != nil {
		return
	}

	if !conn.Ok {
		err = fmt.Errorf("ERR CONN NEGO: %+v", conn)
		return
	}
	log.Printf("-i- CONN ..OK")

	return
}

// Get pulls a message out of the RTM queue and returns it as a Message
func (conn Conn) Get() (m Message, err error) {
	err = websocket.JSON.Receive(conn.Sock, &m)
	if err != nil {
		return
	}
	return
}

// ReplyTo sends a Reply to a Message via Conn.Sock
func (conn Conn) ReplyTo(m *Message, r Reply) (err error) {
	if r.Text == "" {
		err = fmt.Errorf("ERR CONN RPLY: empty Reply.Text")
		return
	}

	r.ID = atomic.AddUint64(&counter, 1)
	r.Channel = m.Channel // TODO is this legit?
	if r.Type == "" {
		r.Type = "message"
	}
	if m.User != "" && r.ReplyToUser {
		r.Text = "<@" + m.User + "> " + r.Text
	}

	log.Printf("<<< [%s] %s\n", r.Type, r.Text)
	return websocket.JSON.Send(conn.Sock, &r)
}

// ReplyWithErr sends a Reply to a Message and logs an error
func (conn Conn) ReplyWithErr(m *Message, msg string, err error) {
	log.Printf("ERR CMD. .ERR: %s", err)
	r := Reply{}
	r.Text = msg
	if r.Text == "" {
		r.Text = "problema!"
	}
	conn.ReplyTo(m, r)
}
