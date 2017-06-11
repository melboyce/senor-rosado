package main

import (
	"math/rand"
	"time"

	"github.com/weirdtales/senor-rosado/slack"
)

var greetings = []string{
	"hola",
	"muy buenos",
	"¿qué tal?",
	"¿qué pasa?",
	"¿qué hubo?",
	"bienvenidos",
	"¿aló?",
	"¿dónde has estado?",
	"¡hace tiempo que no te veo!",
}

// Register ...
func Register() (r string, h string) {
	r = "^(hello|hi)"
	h = "`hello|hi` say hello"
	return
}

// Respond ...
func Respond(m slack.Message, c slack.Conn, matches []string) {
	defer slack.PanicSuppress()
	reply := slack.Reply{}
	reply.Channel = m.Channel
	rand.Seed(time.Now().Unix())
	reply.Text = greetings[rand.Intn(len(greetings))]
	c.Send(m, reply)
}
