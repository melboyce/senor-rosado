package main

import "github.com/weirdtales/senor-rosado/slack"

// Register ...
func Register() string {
	return "^(hello|hi)"
}

// Respond ...
func Respond(m slack.Message, c slack.Conn, matches []string) {
	reply := slack.Reply{}
	reply.Channel = m.Channel
	reply.Text = "hola"
	c.Send(m, reply)
}
