package slack

import (
	"fmt"
	"log"
	"net/url"

	"golang.org/x/net/websocket"
)

// Conn ...
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

// Connect ...
func Connect(token string) (conn Conn, err error) {
	conn.token = url.QueryEscape(token)
	u := fmt.Sprintf(slackURL, conn.token)
	if err = connect(u, &conn); err != nil {
		return
	}
	err = attachSock(&conn)
	return
}

// Get ...
func (conn Conn) Get() (m Message, err error) {
	if conn.Sock == nil {
		panic("!!! CONN SOCK MISSING")
	}
	err = websocket.JSON.Receive(conn.Sock, &m)
	if m.Type == "message" {
		m.UserDetail, err = GetUser(conn.token, m.User)
	}
	return
}

// Send ...
func (conn Conn) Send(reply Reply) (err error) {
	if reply.Text == "" {
		err = fmt.Errorf("Reply.Text is empty")
		return
	}
	log.Printf("<<< [%s:%d] %s", reply.Type, reply.ID, reply.Text)
	return websocket.JSON.Send(conn.Sock, &reply)
}

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

func attachSock(conn *Conn) (err error) {
	log.Printf("-i- CONN SOCK: %s", conn.URL)
	origin := "https://api.slack.com/"
	conn.Sock, err = websocket.Dial(conn.URL, "", origin)
	return
}
