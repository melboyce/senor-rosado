package cmdhello

import "github.com/weirdtales/senor-rosado/slack"

// Register ...
func Register() (re string, help string) {
	re = `^(hello|hi)`
	help = "`hello` say hello"
	return
}

// Respond ...
func Respond(r slack.Reply, c chan slack.Reply) {
	r.Text = "hola!"
	c <- r
}
