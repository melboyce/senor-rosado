package main

import "log"
import "fmt"
import "os"

import "github.com/weirdtales/senor-rosado/slack"


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

    os.Exit(loop(conn))
}

func loop(conn slack.Conn) int {
    reply := slack.Message{}
    replyToUser := true
    for {
        msg, err := conn.Get()
        if err != nil {
            log.Fatal(err)
        }

        if os.Getenv("DEBUG") == "1" {
            log.Printf("[d] msg=%+v", msg)
        }
        if ! msg.Respond {
            continue
        }
        log.Printf(">>> %s %s: %s", msg.Channel, msg.User, msg.Text)

        switch msg.Command {
            case "help", "?", "usage", "-h", "--help":
                reply.Text = "not yet"
            case "toggle-reply-target":
                log.Printf("-?- toggling replyToUser: %v", !replyToUser)
                replyToUser = !replyToUser
        }

        if reply.Text != "" {
            if replyToUser {
                reply.Text = "<@" + msg.User + "> " + reply.Text
            }

            log.Printf("<<< %s %s", msg.Channel, reply.Text)
            err = conn.Send(reply, msg.Channel)
            if err != nil {
                log.Fatal(err)
            }
        }

        reply = slack.Message{}
    }
    return 0
}
