package slack

import (
	"fmt"

	"golang.org/x/net/websocket"
)

var connectURL = "https://slack.com/api/rtm.connect?token=%s"
var wsOrigin = "https://api.slack.com/"

type apiConnect struct {
	OK   bool   `json:"ok"`
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
}

// Conn is a Slack connection
type Conn struct {
	self User
	sock *websocket.Conn
}

// NewConn returns a Conn
func NewConn() (conn Conn) {
	var err error
	var raw apiConnect
	loc := fmt.Sprintf(connectURL, token)

	// hit the main API endpoint
	if err = getJSON(loc, &raw); err != nil {
		panic(err)
	}

	// attach a websocket to send messages on
	if conn.sock, err = dialMessageServer(raw.URL); err != nil {
		panic(err)
	}

	// attach the self User
	if conn.self, err = NewUser(raw.Self.ID); err != nil {
		panic(err)
	}
	return
}

func dialMessageServer(loc string) (sock *websocket.Conn, err error) {
	return websocket.Dial(loc, "", wsOrigin)
}
