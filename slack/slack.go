// Package slack provides an interface for working with slack.com
package slack

import (
        "fmt"
        "strings"

        "encoding/json"
        "io/ioutil"
        "net/http"
        "sync/atomic"

        "golang.org/x/net/websocket"
)


// A Conn represents a slack connection.
type Conn struct {
    Ok    bool     `json:"ok"`
    URL   string   `json:"url"`
    Team  connTeam `json:"team"`
    Self  connSelf `json:"self"`
    Error string   `json:"error"`

    Sock  *websocket.Conn
}

type connTeam struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type connSelf struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// A Message is a slack RTM message object with some meta.
type Message struct {
    ID      uint64 `json:"id"`
    Type    string `json:"type"`
    Subtype string `json:"subtype"`
    Channel string `json:"channel"`
    Text    string `json:"text"`
    User    string `json:"user"`

    Respond    bool
    Target     string
    Command    string
    Subcommand string
    Tail       string
}

// A Reply is another name for a Message
type Reply Message

// see: Conn.Send()
var counter uint64

// Connect to slack and return a useful struct or an error.
func Connect(token string) (slack Conn, err error) {
    url := fmt.Sprintf("https://slack.com/api/rtm.connect?token=%s", token)
    r, err := http.Get(url)
    if err != nil {
        return
    }
    if r.StatusCode != 200 {
        err = fmt.Errorf("Error: status code: %d", r.StatusCode)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        return
    }

    err = json.Unmarshal(body, &slack)
    if err != nil {
        return
    }
    if ! slack.Ok {
        err = fmt.Errorf("Error: slack: %s", slack.Error)
        return
    }

    ws, err := websocket.Dial(slack.URL, "", "https://api.slack.com/")
    if err != nil {
        return
    }
    slack.Sock = ws

    return
}

// Get pulls a message out of the RTM queue and returns it as a Message struct.
func (s Conn) Get() (m Message, err error) {
    err = websocket.JSON.Receive(s.Sock, &m)
    if err != nil {
        return
    }

    if m.Type != "message" {
        return
    }

    // candidate for generalization if this block gets fat (parseMessage)
    words := strings.Split(m.Text, " ")
    if words[0] == "<@" + s.Self.ID + ">" {
        // slack.ChatLoop uses this to ignore messages not targetted at the bot
        m.Respond = true
    }

    n := len(words)

    // these are conveniences used when authoring commands
    m.Target = words[0] // the user that was targetted
    if n > 1 { m.Command    = words[1] }
    if n > 2 { m.Subcommand = words[2] }
    if n > 3 { m.Tail       = strings.Join(words[3:], " ") }

    return
}

// Send pushes a Reply struct into the RTM queue.
func (s Conn) Send(r *Reply, channel string) error {
    r.ID = atomic.AddUint64(&counter, 1)
    r.Channel = channel
    r.Type = "message" // TODO this will bite me later
    return websocket.JSON.Send(s.Sock, &r)
}
