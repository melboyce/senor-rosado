package main

import (
	"log"
	"os"

	"github.com/weirdtales/senor-rosado/slack"
)

func main() {
	// TODO: signals
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Fatal("!!! Need SLACK_TOKEN in the env, man.\n")
		os.Exit(1)
	}

	conn, err := slack.Connect(token)
	if err != nil {
		log.Printf("%+v\n", conn)
		log.Fatal(err)
	}

	log.Println("-i- starting up...")
	slack.ChatLoop(conn) // TODO ChatLoop doesn't return
}
