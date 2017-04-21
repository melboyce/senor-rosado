package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "sync/atomic"

        "golang.org/x/net/websocket"
)


// root response from: https://slack.com/api/rtm.connect
type respRtmConnect struct {
    Ok    bool        `json:"ok"`
    URL   string      `json:"url"`
    Team  respRtmTeam `json:"team"`
    Self  respRtmSelf `json:"self"`
    Error string      `json:"error"`
}

// team struct for respRtmConnect
type respRtmTeam struct {
    ID string `json:"id"`
}

// self struct for respRtmConnect
type respRtmSelf struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// Message structure
type Message struct {
    ID      uint64 `json:"id"`
    Type    string `json:"type"`
    Channel string `json:"channel"`
    Text    string `json:"text"`
}

// used to generate IDs
var counter uint64

// returns a websocket URL and the user ID
func connect(token string) (wsurl, id string, err error) {
    url := fmt.Sprintf("https://slack.com/api/rtm.connect?token=%s", token)
    resp, err := http.Get(url)
    if err != nil {
        return
    }
    if resp.StatusCode != 200 {
        err = fmt.Errorf("Error (status): %d", resp.StatusCode)
        return
    }
    body, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        return
    }

    var rtmResp respRtmConnect
    err = json.Unmarshal(body, &rtmResp)
    if err != nil {
        return
    }
    if ! rtmResp.Ok {
        err = fmt.Errorf("Error (slack): %s", rtmResp.Error)
        return
    }

    return rtmResp.URL, rtmResp.Self.ID, nil
}

// returns a websocket connection and the user ID
func api(token string) (*websocket.Conn, string) {
    wsurl, id, err := connect(token)
    if err != nil {
        log.Fatal(err)
    }
    ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
    if err != nil {
        log.Fatal(err)
    }

    return ws, id
}

// returns a message
func get(ws *websocket.Conn) (m Message, err error) {
    err = websocket.JSON.Receive(ws, &m)
    return m, err
}

// sends a message
func post(ws *websocket.Conn, m Message) error {
    m.ID = atomic.AddUint64(&counter, 1)
    return websocket.JSON.Send(ws, m)
}
