package slack

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUserGetFromCache ...
func TestUserGetFromCache(t *testing.T) {
	expected := User{"abc", "Alpha Bravo Charlie"}
	setCachedUser(expected)
	actual, ok := getCachedUser("abc")

	if !ok {
		t.Fatalf("Expected ok == true; not the case")
	}

	if actual != expected {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

// TestUserGet ...
func TestUserGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{
    "ok": true,
    "user": {
		"id": "U023BECGF",
		"name": "bobby",
		"deleted": false,
		"color": "9f69e7",
		"profile": {
			"avatar_hash": "ge3b51ca72de",
			"status_emoji": ":mountain_railway:",
			"status_text": "riding a train",
			"first_name": "Bobby",
			"last_name": "Tables",
			"real_name": "Bobby Tables",
			"email": "bobby@slack.com",
			"skype": "my-skype-name",
			"phone": "+1 (123) 456 7890",
			"image_24": "https:\/\/...",
			"image_32": "https:\/\/...",
			"image_48": "https:\/\/...",
			"image_72": "https:\/\/...",
			"image_192": "https:\/\/...",
			"image_512": "https:\/\/..."
		},
		"is_admin": true,
		"is_owner": true,
		"is_primary_owner": true,
		"is_restricted": false,
		"is_ultra_restricted": false,
		"updated": 1490054400,
		"has_2fa": false,
		"two_factor_type": "sms"
    }
}`
		if r.URL.Query().Get("user") == "U023BECGF" {
			fmt.Fprintln(w, resp)
		} else {
			fmt.Fprintln(w, `{"ok":false,"user":{}}`)
		}
	}))
	defer ts.Close()

	userURL = ts.URL + "?token=%s&user=%s"
	actual, err := GetUser("U023BECGF")
	if err != nil {
		t.Fatalf("err: %+v", err)
	}
	if actual.name != "bobby" {
		t.Fatalf("Expected 'bobby', got '%s'", actual.name)
	}
	if actual.name != actual.String() {
		t.Fatalf("Expected User.String() == User.name (Stringer test), but no.")
	}

	actual, err = GetUser("xxx")
	if err != nil {
		t.Fatalf("err: %+v", err)
	}
	if actual.name != "" {
		t.Fatalf("Expected '', got '%s'", actual.name)
	}

	// getUsers(ids []string)
	ids := []string{"U023BECGF", "xxx"}
	actualUsers, err := getUsers(ids)
	if err != nil {
		t.Fatalf("err: %+v", err)
	}
	if len(actualUsers) != 2 {
		t.Fatalf("Expected 2 results from getUsers(), got %d", len(actualUsers))
	}
}
