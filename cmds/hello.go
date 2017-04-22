package cmds

import "github.com/weirdtales/senor-rosado/slack"


// Hello says hello
func Hello(m slack.Message, r *slack.Reply) error {
    r.Text = "Hola"
    return nil
}
