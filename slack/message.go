package slack

// Message ...
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`

	Respond bool
	Cmd     string
	Args    []string
}

// Reply ...
type Reply struct {
	*Message
}

var counter uint64
