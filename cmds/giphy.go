package cmds

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/weirdtales/senor-rosado/slack"
)

// GiphyRegister for giphy cartridge
func GiphyRegister() (r string, h string) {
	r = "^(giphy) (.+)"
	h = "`giphy kristin bell sloth` lucky draw gif from giphy"
	return
}

type giphySearch struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

// GiphyRespond finds a Gif on Giphy
func GiphyRespond(c *slack.Conn, m *slack.Message) {
	token := os.Getenv("GIPHY_TOKEN")
	if token == "" {
		c.ReplyTo(m, slack.Reply{Text: "No token in env."})
		return
	}

	q := url.QueryEscape(strings.Join(m.Args, " "))
	u := "http://api.giphy.com/v1/gifs/search?q=%s&api_key=%s"
	u = fmt.Sprintf(u, q, token)

	var gifs giphySearch
	if err := slack.GetJSON(u, &gifs); err != nil {
		c.ReplyWithErr(m, "problem loading response from giphy", err)
		return
	}

	r := slack.Reply{Text: "hoy no"}
	if len(gifs.Data) > 0 {
		rand.Seed(time.Now().Unix())
		maxRand := 3
		if len(gifs.Data) < maxRand {
			maxRand = len(gifs.Data)
		}
		r.Text = gifs.Data[rand.Intn(maxRand)].URL
	}
	c.ReplyTo(m, r)
}
