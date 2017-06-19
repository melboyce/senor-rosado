package slack

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}
var slackURL = "https://slack.com/api/rtm.connect?token=%s"

// Message is an incoming Slack Message
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	SubType string `json:"subtype"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`

	Cmd         string
	Args        []string
	ReplyToUser bool
	Respond     bool
	Target      string
	Matches     [][]string
}

// Reply is an outgoing Message
type Reply Message

// Command is a bot command - see package: cmds
type Command struct {
	Re       *regexp.Regexp
	Help     string
	Register func() (string, string)
	Respond  func(*Conn, *Message)
}

// Init loads internal and any configured external commands
func Init(conn *Conn, cmds []Command) (err error) {
	var i int
	attachSock(conn)

	var patt string
	for _, cmd := range getBuiltins() {
		patt, cmd.Help = cmd.Register()
		cmd.Re = regexp.MustCompile(patt)
		conn.Commands = append(conn.Commands, cmd)
		i++
	}
	for _, cmd := range cmds {
		patt, cmd.Help = cmd.Register()
		cmd.Re = regexp.MustCompile(patt)
		conn.Commands = append(conn.Commands, cmd)
		i++
	}

	log.Printf("-i- INIT ..OK: %d command(s) registered", i)
	return
}

// Loop runs a hard loop against the RTM queue and processes the input
func Loop(conn *Conn) (err error) {
	log.Printf("-i- LOOP STRT")
	for {
		m, err := conn.Get()
		if err != nil {
			panic(err)
		}

		if m.Type != "message" {
			continue
		}

		log.Printf(">>> [%s] %s", m.Type, m.Text)

		if len(conn.Commands) < 1 {
			continue
		}

		parseMessage(conn, &m)

		if !m.Respond {
			continue
		}

		matchCommands(conn, m)
	}
}

func matchCommands(conn *Conn, m Message) {
	for _, cmd := range conn.Commands {
		m.Matches = cmd.Re.FindAllStringSubmatch(m.Text, -1)
		if len(m.Matches) > 0 {
			if os.Getenv("DEBUG") == "1" {
				log.Printf("-d- RGXP ..OK: %s =~ %s", cmd.Re, m.Text)
			}
			go cmd.Respond(conn, &m)
		}
	}
}

func getBuiltins() (cmds []Command) {
	return []Command{
		Command{
			Register: BuiltinHelpRegister,
			Respond:  BuiltinHelpRespond,
		},
	}
}

func attachSock(conn *Conn) {
	var err error
	conn.Origin = "https://api.slack.com/"
	conn.Sock, err = websocket.Dial(conn.URL, "", conn.Origin)
	if err != nil {
		panic(err) // TODO something more graceful
	}
}

func parseMessage(conn *Conn, m *Message) {
	if m.Type != "message" || m.Text == "" {
		return
	}

	words := strings.Split(m.Text, " ")

	// msg directed at @bot
	if words[0] == "<@"+conn.Self.ID+">" {
		m.Respond = true
		m.Target = words[0]
		m.Text = strings.Join(words[1:], " ")
	}

	// msg is private
	if strings.HasPrefix(m.Channel, "D") {
		m.Respond = true
		m.Target = conn.Self.ID
		words = append([]string{conn.Self.ID}, words...)
	}

	n := len(words)
	if n > 1 {
		m.Cmd = words[1]
		m.Args = words[1:]
	}

	if os.Getenv("DEBUG") == "1" {
		log.Printf("-d- CORE PARS: %+v", m)
	}

	return
}
