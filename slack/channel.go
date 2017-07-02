package slack

import (
	"fmt"
)

var channelURL = "https://slack.com/api/channels.info?token=%s&channel=%s"

type apiChannel struct {
	Channel struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Members []string `json:"members"`
		Topic   struct {
			Value   string `json:"value"`
			Creator string `json:"creator"`
			LastSet int    `json:"last_set"`
		} `json:"topic"`
		Purpose struct {
			Value   string `json:"value"`
			Creator string `json:"creator"`
			LastSet int    `json:"last_set"`
		} `json:"purpose"`
	} `json:"channel"`
}

// Channel is a Slack channel representation.
type Channel struct {
	id      string
	name    string
	topic   string
	purpose string
	members []User
}

// TODO cache expiry, lock
var channelCache = make(map[string]Channel)

// GetChannel returns a possibly cached Channel.
func GetChannel(id string) (c Channel, err error) {
	c, ok := channelCache[id]
	if ok {
		return
	}
	loc := fmt.Sprintf(channelURL, token, id)
	var raw apiChannel
	if err = getJSON(loc, &raw); err != nil {
		panic(err)
	}

	c = Channel{
		id:      raw.Channel.ID,
		name:    raw.Channel.Name,
		topic:   raw.Channel.Topic.Value,
		purpose: raw.Channel.Purpose.Value,
	}
	c.members, err = getUsers(raw.Channel.Members)
	channelCache[id] = c
	return
}

// String ...
func (c Channel) String() string {
	return c.name
}
