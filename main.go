package main

import (
	"log"
	"os"

	"github.com/weirdtales/senor-rosado/cmds"
	"github.com/weirdtales/senor-rosado/cmdweather"
	"github.com/weirdtales/senor-rosado/slack"
)

var commands = []slack.Command{
	slack.Command{
		Register: cmds.HelloRegister,
		Respond:  cmds.HelloRespond,
	},
	slack.Command{
		Register: cmds.GiphyRegister,
		Respond:  cmds.GiphyRespond,
	},
	slack.Command{
		Register: cmds.FixerRegister,
		Respond:  cmds.FixerRespond,
	},
	slack.Command{
		Register: cmdweather.Register,
		Respond:  cmdweather.Respond,
	},
}

func main() {
	// TODO: signals
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Printf("!!! missing SLACK_TOKEN env var\n")
		os.Exit(127)
	}

	conn, err := slack.Connect(token)
	if err != nil {
		log.Printf("%+v\n", conn)
		log.Print(err)
		os.Exit(126)
	}

	slack.Init(&conn, commands)
	slack.Loop(&conn)
}
