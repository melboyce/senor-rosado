package cmds

import (
	"math/rand"
	"time"

	"github.com/weirdtales/senor-rosado/slack"
)

var helloGreetings = []string{
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

// HelloRegister ...
func HelloRegister() (re string, help string) {
	re = `^(hello|hi|hola|ola|hey|oi|yo).*`
	help = "`hello` say hello!"
	return
}

// HelloRespond ...
func HelloRespond(conn *slack.Conn, m *slack.Message) {
	r := slack.Reply{}
	rand.Seed(time.Now().Unix())
	r.Text = helloGreetings[rand.Intn(len(helloGreetings))]
	conn.ReplyTo(m, r)
}
