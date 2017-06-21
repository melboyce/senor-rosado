package cmdhello

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
func Register() (re string, help string) {
	re = `^(hello|hi|hola|ola|hey|hi|yo|oi)`
	help = "`hello` say hello"
	return
}

// Respond ...
func Respond(r slack.Reply, c chan slack.Reply) {
	rand.Seed(time.Now().Unix())
	r.Text = greetings[rand.Intn(len(greetings))]
	c <- r
}
