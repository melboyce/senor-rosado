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

	Sock *websocket.Conn
}

var slackURL = "https://slack.com/api/rtm.connect?token=%s"

// Connect ...
func Connect(token string) (conn Conn, err error) {
	u := fmt.Sprintf(slackURL, url.QueryEscape(token))
	conn, err = connect(u)
	if err != nil {
		return
	}
	err = attachSock(conn)
	return
}

// Get ...
func (conn Conn) Get() (m Message, err error) {
	err = websocket.JSON.Receive(conn.Sock, &m)
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

func connect(u string) (conn Conn, err error) {
	log.Printf("-i- CONN STRT")
	if err = GetJSON(u, &conn); err != nil {
		return
	}
	if !conn.OK {
		err = fmt.Errorf("Conn.OK is false: %+v", conn)
	}
	return
}

func attachSock(conn Conn) (err error) {
	log.Printf("-i- CONN SOCK: %s", conn.URL)
	origin := "https://api.slack.com/"
	conn.Sock, err = websocket.Dial(conn.URL, "", origin)
	return
}
