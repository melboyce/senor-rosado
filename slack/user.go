package slack

import "fmt"

var userURL = "https://slack.com/api/users.info?token=%s&user=%s"

type apiUser struct {
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
			Tz            string `json:"tz"`
			TzLabel       string `json:"tz_label"`
			TzOffset      int    `json:"tz_offset"`
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
		Has2Fa  bool `json:"has_2fa"`
	} `json:"user"`
}

// User is... well, a user.
type User struct {
	id   string
	name string
}

// TODO cache expiry, lock
var userCache = make(map[string]User)

// GetUser returns a possibly cached User.
func GetUser(id string) (u User, err error) {
	u, ok := userCache[id]
	if ok {
		return
	}
	loc := fmt.Sprintf(userURL, token, id)
	var raw apiUser
	if err = getJSON(loc, &raw); err != nil {
		panic(err)
	}
	u = User{
		id:   raw.User.ID,
		name: raw.User.Name,
	}
	userCache[id] = u
	return
}

// getUsers takes a slice of user IDs (strings) and returns a
// slice of User structs.
func getUsers(ids []string) (u []User, err error) {
	var n User
	for _, id := range ids {
		if n, err = GetUser(id); err != nil {
			continue
		}
		u = append(u, n)
	}
	return
}

// String ...
func (u User) String() string {
	return u.name
}
