package slack

import (
	"golang.org/x/net/websocket"
)

type apiMessage struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Text    string `json:"text"`
}

// Message is an incoming message.
type Message struct {
	apiMessage
	SChannel Channel
	SUser    User
}

// getMessages pushes Slack Messages into `out`.
func getMessages(conn Conn, out chan Message, quit chan int) {
	var (
		m   Message
		err error
	)
	for {
		m = Message{}
		if err = websocket.JSON.Receive(conn.sock, &m.apiMessage); err != nil {
			quit <- 1
			break
		}
		// TODO error handling
		if m.Type == "message" {
			m.SChannel, _ = NewChannel(m.Channel)
			m.SUser, _ = NewUser(m.User)
		}
		out <- m
	}
}
