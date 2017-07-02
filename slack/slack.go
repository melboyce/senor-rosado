package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}
var token string

// Start the Slack session.
func Start(t string) int {
	token = t
	var (
		exit = 0
		conn = NewConn()
		quit = make(chan int)
		in   = make(chan Message)
		m    Message
	)
	go getMessages(conn.sock, in, quit)

Loop:
	for {
		select {
		case m = <-in:
			if m.Type != "message" || m.Text == "" {
				continue
			}
			processMessage(conn, &m)
		case exit = <-quit:
			break Loop
		}
	}

	return exit
}

// getJSON marshals the JSON response from `loc` into `&t`.
func getJSON(loc string, t interface{}) (err error) {
	u, err := url.Parse(loc)
	if err != nil {
		return
	}

	log.Printf("JSON %s://%s%s", u.Scheme, u.Host, u.Path)
	r, err := httpClient.Get(u.String())
	if err != nil {
		return
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		return
	}

	if r.StatusCode != 200 {
		err = fmt.Errorf("getJSON: status != 200")
		return
	}
	return
}
