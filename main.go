package main

import (
	"log"
	"os"

	"github.com/weirdtales/senor-rosado/slack"
)

func main() {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		log.Fatalf("FATL no token in env")
	}
	os.Exit(slack.Start(token))
}
