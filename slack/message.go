package slack

import (
	"log"
	"strings"

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
	Obj struct {
		Channel Channel
		User    User
	}
	Targeted bool
	Cmd      string
	Args     []string
}

// Reply is an outgoing message
type Reply apiMessage

// getMessages takes messages from `sock` and pushes them to `out`.
func getMessages(sock *websocket.Conn, out chan Message, quit chan int) {
	var (
		m   Message
		err error
	)
	for {
		m = Message{}
		if err = websocket.JSON.Receive(sock, &m.apiMessage); err != nil {
			quit <- 1
			break
		}

		if m.Type == "message" {
			m.Obj.Channel, _ = GetChannel(m.Channel)
			m.Obj.User, _ = GetUser(m.User)
			log.Printf(" >>> #%s @%s: %s", m.Obj.Channel, m.Obj.User, m.Text)
		}

		out <- m
	}
}

func processMessage(conn Conn, m *Message) {
	words := strings.Split(m.Text, " ")
	if len(words) > -1 && words[0] == "<@"+conn.self.id+">" {
		m.Targeted = true
		words = words[1:]
	}
	if len(words) > 0 {
		m.Cmd = words[0]
	}
	if len(words) > 1 {
		m.Args = words[1:]
	}
}
