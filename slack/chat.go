package slack

// TODO share logger with caller
import "log"
import "os"


// ChatFn is used to make a map of command functions
type ChatFn func(Message, *Reply) error

// ChatLoop enters a hard loop that reads off messages and processes them.
// Note: this function doesn't return errors.
func ChatLoop(conn Conn, cmds map[string]ChatFn) int {
    var reply Reply
    replyToUser := true

    for {
        reply = Reply{}
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

        // commands handled by the chat loop
        // this is for commands that need to impact the loop directly (options, etc)
        switch msg.Command {
            case "help", "?", "usage", "-h", "--help":
                // TODO see if there's some way to reflect the commands to make help programmatic
                reply.Text = "not yet"
            case "opt-replytouser", "opt-rtu":
                log.Printf("-?- toggling replyToUser: %v", !replyToUser)
                replyToUser = !replyToUser
        }

        // check if command is in the map
        if cmd, ok := cmds[msg.Command]; ok {
            err = cmd(msg, &reply)
        }
        if err != nil {
            log.Fatal(err)
        }

        if reply.Text == "" {
            continue
        }

        if replyToUser {
            reply.Text = "<@" + msg.User + "> " + reply.Text
        }

        log.Printf("<<< %s %s", msg.Channel, reply.Text)
        err = conn.Send(&reply, msg.Channel)
        if err != nil {
            log.Fatal(err)
        }
    }
    return 0 // TODO not currently reachable
}
