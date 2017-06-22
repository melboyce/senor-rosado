package slack

import "fmt"

// User is a response from the Slack user.info endpoint.
type User struct {
	Ok   bool `json:"ok"`
	User struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Deleted bool   `json:"deleted"`
		Color   string `json:"color"`
		Profile struct {
			AvatarHash    string `json:"avatar_hash"`
			CurrentStatus string `json:"current_status"`
			FirstName     string `json:"first_name"`
			LastName      string `json:"last_name"`
			RealName      string `json:"real_name"`
			Email         string `json:"email"`
			Skype         string `json:"skype"`
			Phone         string `json:"phone"`
			Image24       string `json:"image_24"`
			Image32       string `json:"image_32"`
			Image48       string `json:"image_48"`
			Image72       string `json:"image_72"`
			Image192      string `json:"image_192"`
		} `json:"profile"`
		IsAdmin bool `json:"is_admin"`
		IsOwner bool `json:"is_owner"`
		Updated int  `json:"updated"`
		Has2Fa  bool `json:"has_2fa"`
	} `json:"user"`
}

// UserLookup memoizes user.info API hits.
var UserLookup = map[string]*User{}

// TODO programatic URL construction?
var slackUserURL = "https://slack.com/api/users.info?token=%s&user=%s"

// GetUser returns a pointer to a User.
func GetUser(token string, id string) (user *User, err error) {
	if user, ok := UserLookup[id]; ok {
		return user, nil
	}
	u := fmt.Sprintf(slackUserURL, token, id)
	if err = GetJSON(u, &user); err != nil {
		return
	}
	UserLookup[id] = user
	return
}
