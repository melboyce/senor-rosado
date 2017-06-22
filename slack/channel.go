package slack

import "fmt"

// Channel represents a channel.info response from the Slack API
type Channel struct {
	Ok      bool `json:"ok"`
	Channel struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		IsChannel      bool   `json:"is_channel"`
		Created        int    `json:"created"`
		Creator        string `json:"creator"`
		IsArchived     bool   `json:"is_archived"`
		IsGeneral      bool   `json:"is_general"`
		NameNormalized string `json:"name_normalized"`
		IsShared       bool   `json:"is_shared"`
		IsOrgShared    bool   `json:"is_org_shared"`
		IsMember       bool   `json:"is_member"`
		LastRead       string `json:"last_read"`
		Latest         struct {
			Type string `json:"type"`
			User string `json:"user"`
			Text string `json:"text"`
			Ts   string `json:"ts"`
		} `json:"latest"`
		UnreadCount        int      `json:"unread_count"`
		UnreadCountDisplay int      `json:"unread_count_display"`
		Members            []string `json:"members"`
		Topic              struct {
			Value   string `json:"value"`
			Creator string `json:"creator"`
			LastSet int    `json:"last_set"`
		} `json:"topic"`
		Purpose struct {
			Value   string `json:"value"`
			Creator string `json:"creator"`
			LastSet int    `json:"last_set"`
		} `json:"purpose"`
		PreviousNames []interface{} `json:"previous_names"`
	} `json:"channel"`
}

// ChannelLookup memoizes channel.info calls
var ChannelLookup = map[string]*Channel{}

var slackChannelURL = "https://slack.com/api/channels.info?token=%s&channel=%s"

// GetChannel hits Slack's channel.info endpoint, potentially stores the
// result in ChannelLookup, and returns a pointer to a Channel struct.
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
