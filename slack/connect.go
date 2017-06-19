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
}

// Connect ...
func Connect(token string) (conn Conn, err error) {
	log.Printf("-i- CONN STRT")
	u := fmt.Sprintf(slackURL, url.QueryEscape(token))
	if err = GetJSON(u, &conn); err != nil {
		return
	}

	if !conn.OK {
		err = fmt.Errorf("/!\\ CONN NEGO: %+v", conn)
		return
	}

	log.Printf("-i- CONN SOCK: %s", conn.URL)
	origin := "https://api.slack.com/"
	conn.Sock, err = websocket.Dial(conn.URL, "", origin)
	if err != nil {
		return
	}

	log.Printf("-i- CONN OKAY")
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
		err = fmt.Errorf("/!\\ CONN SEND: empty response (Reply.Text)")
		return
	}
	log.Printf("<<< [%s:%d] %s", reply.Type, reply.ID, reply.Text)
	return websocket.JSON.Send(conn.Sock, &r)
}
