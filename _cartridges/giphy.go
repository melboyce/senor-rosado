package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/weirdtales/senor-rosado/slack"
)

// Register for giphy cartridge
func Register() (r string, h string) {
	r = "^giphy "
	h = "`giphy kristin bell sloth` lucky draw gif from giphy"
	return
}

type giphySearch struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

// Respond finds a Gif on Giphy
func Respond(m slack.Message, c slack.Conn, matches []string) {
	reply := slack.Reply{}
	reply.Channel = m.Channel

	token := os.Getenv("GIPHY_TOKEN")
	if token == "" {
		reply.Text = "Need a Giphy token in the env pls."
		c.Send(m, reply)
		return
	}

	url := "http://api.giphy.com/v1/gifs/search?q=%s&api_key=%s"
	url = fmt.Sprintf(url, strings.Replace(m.Tail, " ", "+", -1), token)

	var gifs giphySearch
	err := slack.GetJSON(url, &gifs)
	if err != nil {
		log.Printf("ERR %s\n", err)
		reply.Text = "problema!"
		c.Send(m, reply)
		return
	}

	reply.Text = "hoy no"
	if len(gifs.Data) > 0 {
		rand.Seed(time.Now().Unix())
		maxRand := 3
		if len(gifs.Data) < maxRand {
			maxRand = len(gifs.Data)
		}
		reply.Text = gifs.Data[rand.Intn(maxRand)].URL
	}
	c.Send(m, reply)
}
