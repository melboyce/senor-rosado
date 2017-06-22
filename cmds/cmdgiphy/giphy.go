package cmdgiphy

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/weirdtales/senor-rosado/slack"
)

// Register ...
func Register() (r string, h string) {
	r = `^(giphy) (.+)`
	h = "`giphy kristin bell sloth` lucky draw gif from giphy"
	return
}

type giphySearch struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

// Respond finds a Gif on Giphy
func Respond(r slack.Reply, c chan slack.Reply) {
	token := os.Getenv("GIPHY_TOKEN")
	if token == "" {
		r.Text = "no Giphy token"
		c <- r
		return
	}

	q := url.QueryEscape(strings.Join(r.Args, " "))
	u := "http://api.giphy.com/v1/gifs/search?q=%s&api_key=%s"
	u = fmt.Sprintf(u, q, token)

	var gifs giphySearch
	if err := slack.GetJSON(u, &gifs); err != nil {
		r.Text = "problem getting response from giphy"
		c <- r
		return
	}

	r.Text = "hoy no"
	if len(gifs.Data) > 0 {
		rand.Seed(time.Now().Unix())
		maxRand := 3
		if len(gifs.Data) < maxRand {
			maxRand = len(gifs.Data)
		}
		r.Text = gifs.Data[rand.Intn(maxRand)].URL
	}
	c <- r
}
