package main

import (
	"log"
	"os"

	"github.com/weirdtales/senor-rosado/cmds/cmdhello"
	"github.com/weirdtales/senor-rosado/slack"
)

// TODO plugins
var commands = []slack.Command{
	slack.Command{
		Register: cmdhello.Register,
		Respond:  cmdhello.Respond,
	},
}

func main() {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Printf("!!! MAIN MAIN: missing SLACK_TOKEN")
		os.Exit(127)
	}

	conn, err := slack.Connect(token)
	if err != nil {
		log.Printf("!!! MAIN CONN: %s", err)
		os.Exit(126)
	}

	slack.Loop(&conn, commands)
}
