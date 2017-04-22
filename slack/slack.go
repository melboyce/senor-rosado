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


// Conn represents a slack connection
// root response from: https://slack.com/api/rtm.connect
type Conn struct {
    Ok    bool     `json:"ok"`
    URL   string   `json:"url"`
    Team  connTeam `json:"team"`
    Self  connSelf `json:"self"`
    Error string   `json:"error"`

    Sock  *websocket.Conn
}

// team struct for slackConn
type connTeam struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// self struct for slackConn
type connSelf struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// Message structure
type Message struct {
    ID      uint64 `json:"id"`
    Type    string `json:"type"`
    Subtype string `json:"subtype"`
    Channel string `json:"channel"`
    Text    string `json:"text"`
    User    string `json:"user"`

    Respond bool
    Command string
}

// used to generate IDs
var counter uint64

// Connect to slack and return a useful struct
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
    r.Body.Close()
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

// Get pulls a message out of the RTM queue and returns it as a Message struct
func (s Conn) Get() (m Message, err error) {
    err = websocket.JSON.Receive(s.Sock, &m)
    if err != nil {
        return
    }

    if m.Type == "message" {
        words := strings.Split(m.Text, " ")
        if words[0] == "<@" + s.Self.ID + ">" {
            m.Respond = true
            m.Command = words[1]
        }
    }
    return
}

// Send pushes a Message struct into the RTM queue
func (s Conn) Send(m Message, channel string) error {
    m.ID = atomic.AddUint64(&counter, 1)
    m.Channel = channel
    m.Type = "message"
    return websocket.JSON.Send(s.Sock, m)
}
