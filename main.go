package main

import "log"
import "fmt"
import "os"

import "github.com/weirdtales/senor-rosado/slack"
import "github.com/weirdtales/senor-rosado/cmds"


// TODO goroutines for these
var cmdMap = map[string]slack.ChatFn{
    "weather": func(m slack.Message, r *slack.Reply) error { return cmds.Weather(m, r) },
    "hello":   func(m slack.Message, r *slack.Reply) error { return cmds.Hello(m, r) },
}

func main() {
    // TODO: signals
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "usage: %s <TOKEN>\n", os.Args[0])
        os.Exit(1)
    }

    conn, err := slack.Connect(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    os.Exit(slack.ChatLoop(conn, cmdMap))
}
