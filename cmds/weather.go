package cmds

import "github.com/weirdtales/senor-rosado/slack"


// Weather reports the current forecast for an airport code
func Weather(m slack.Message, r *slack.Reply) error {
    r.Text = "YmmL??"
    return nil
}

// Hello says hello
func Hello(m slack.Message, r *slack.Reply) error {
    r.Text = "Hola"
    return nil
}
