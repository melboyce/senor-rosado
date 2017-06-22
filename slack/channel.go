package slack

import "fmt"

// Channel ...
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ChannelLookup ...
var ChannelLookup = map[string]*Channel{}

var slackChannelURL = "https://slack.com/api/channels.info?token=%s&channel=%s"

// GetChannel ...
func GetChannel(token string, id string) (channel *Channel, err error) {
	if channel, ok := ChannelLookup[id]; ok {
		return channel, nil
	}
	u := fmt.Sprintf(slackChannelURL, token, id)
	if err = GetJSON(u, &channel); err != nil {
		return
	}
	ChannelLookup[id] = channel
	return
}
