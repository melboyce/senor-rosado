// Package slack provides an interface for working with slack.com
package slack

import "log"
import "os"


// A ChatFn is a function for a chat command.
type ChatFn func(Message, Conn)


// ChatLoop enters a hard loop that reads off messages and processes them.
func ChatLoop(conn Conn, cmds map[string]ChatFn) int {
    for {
        msg, err := conn.Get()
        if err != nil {
            log.Fatal(err)
        }

        if os.Getenv("DEBUG") == "1" {
            log.Printf("[d] msg=%+v", msg)
        }
        if ! msg.Respond {
            // Message.Respond: if true, message is targetted at bot
            continue
        }
        log.Printf(">>> %s %s: %s", msg.Channel, msg.User, msg.Text)

        // check if command is in the map
        if cmd, ok := cmds[msg.Command]; ok {
            cmd(msg, conn)
        }
    }
    return 0 // TODO not currently reachable
}
