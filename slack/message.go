package slack

import (
	"strings"
	"sync/atomic"
)

// Message ...
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`

	UserDetail *User
	SelfID     string
	Respond    bool
}

// Reply ...
type Reply struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`

	Message Message
	Cmd     string
	Args    []string
}

var counter uint64

// GetReply ...
func GetReply(m Message) (r Reply) {
	r.ID = atomic.AddUint64(&counter, 1)
	r.Type = "message"
	r.Channel = m.Channel
	r.Message = m

	words := strings.Split(m.Text, " ")
	if len(words) >= 0 {
		r.Cmd = words[0]
	}
	if len(words) > 0 {
		r.Args = words[1:]
	}
	return
}
