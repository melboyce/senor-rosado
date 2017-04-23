package cmds

import "time"

import "github.com/weirdtales/senor-rosado/slack"


// Hello says hello
func Hello(m slack.Message, c slack.Conn) {
    time.Sleep(120 * time.Second)
    reply := slack.Reply{}
    reply.Channel = m.Channel
    reply.Text = "hola"
    c.Send(m, reply)
}
